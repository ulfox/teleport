/*
Copyright 2021 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package alpnproxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	tracessh "github.com/gravitational/teleport/api/observability/tracing/ssh"
	"github.com/gravitational/teleport/lib/client"
	"github.com/gravitational/teleport/lib/srv/alpnproxy/common"
	"github.com/gravitational/teleport/lib/utils"
	"github.com/gravitational/teleport/lib/utils/aws"
)

// LocalProxy allows upgrading incoming connection to TLS where custom TLS values are set SNI ALPN and
// updated connection is forwarded to remote ALPN SNI teleport proxy service.
type LocalProxy struct {
	cfg     LocalProxyConfig
	context context.Context
	cancel  context.CancelFunc
}

// LocalProxyConfig is configuration for LocalProxy.
type LocalProxyConfig struct {
	// RemoteProxyAddr is the downstream destination address of remote ALPN proxy service.
	RemoteProxyAddr string
	// Protocol set for the upstream TLS connection.
	Protocols []common.Protocol
	// InsecureSkipTLSVerify turns off verification for x509 upstream ALPN proxy service certificate.
	InsecureSkipVerify bool
	// Listener is listener running on local machine.
	Listener net.Listener
	// SNI is a ServerName value set for upstream TLS connection.
	SNI string
	// ParentContext is a parent context, used to signal global closure>
	ParentContext context.Context
	// SSHUser is an SSH username.
	SSHUser string
	// SSHUserHost is user host requested by ssh subsystem.
	SSHUserHost string
	// SSHHostKeyCallback is the function type used for verifying server keys.
	SSHHostKeyCallback ssh.HostKeyCallback
	// SSHTrustedCluster allows selecting trusted cluster ssh subsystem request.
	SSHTrustedCluster string
	// ClientTLSConfig is a client TLS configuration used during establishing
	// connection to the RemoteProxyAddr.
	ClientTLSConfig *tls.Config
	// Certs are the client certificates used to connect to the remote Teleport Proxy.
	Certs []tls.Certificate
	// AWSCredentials are AWS Credentials used by LocalProxy for request's signature verification.
	AWSCredentials *credentials.Credentials
}

// CheckAndSetDefaults verifies the constraints for LocalProxyConfig.
func (cfg *LocalProxyConfig) CheckAndSetDefaults() error {
	if cfg.RemoteProxyAddr == "" {
		return trace.BadParameter("missing remote proxy address")
	}
	if len(cfg.Protocols) == 0 {
		return trace.BadParameter("missing protocol")
	}
	if cfg.ParentContext == nil {
		return trace.BadParameter("missing parent context")
	}
	return nil
}

func (cfg *LocalProxyConfig) GetProtocols() []string {
	protos := make([]string, 0, len(cfg.Protocols))

	for _, proto := range cfg.Protocols {
		protos = append(protos, string(proto))
	}

	return protos
}

// NewLocalProxy creates a new instance of LocalProxy.
func NewLocalProxy(cfg LocalProxyConfig) (*LocalProxy, error) {
	if err := cfg.CheckAndSetDefaults(); err != nil {
		return nil, trace.Wrap(err)
	}

	ctx, cancel := context.WithCancel(cfg.ParentContext)
	return &LocalProxy{
		cfg:     cfg,
		context: ctx,
		cancel:  cancel,
	}, nil
}

// SSHProxy is equivalent of `ssh -o 'ForwardAgent yes' -p port  %r@host -s proxy:%h:%p` but established SSH
// connection to RemoteProxyAddr is wrapped with TLS protocol.
func (l *LocalProxy) SSHProxy(ctx context.Context, localAgent *client.LocalKeyAgent) error {
	if l.cfg.ClientTLSConfig == nil {
		return trace.BadParameter("client TLS config is missing")
	}

	clientTLSConfig := l.cfg.ClientTLSConfig.Clone()
	clientTLSConfig.NextProtos = l.cfg.GetProtocols()
	clientTLSConfig.InsecureSkipVerify = l.cfg.InsecureSkipVerify
	clientTLSConfig.ServerName = l.cfg.SNI

	upstreamConn, err := tls.Dial("tcp", l.cfg.RemoteProxyAddr, clientTLSConfig)
	if err != nil {
		return trace.Wrap(err)
	}
	defer upstreamConn.Close()

	client, err := makeSSHClient(ctx, upstreamConn, l.cfg.RemoteProxyAddr, &ssh.ClientConfig{
		User: l.cfg.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(localAgent.Signers),
		},
		HostKeyCallback: l.cfg.SSHHostKeyCallback,
	})
	if err != nil {
		return trace.Wrap(err)
	}
	defer client.Close()

	sess, err := client.NewSession(ctx)
	if err != nil {
		return trace.Wrap(err)
	}
	defer sess.Close()

	err = agent.ForwardToAgent(client.Client, localAgent)
	if err != nil {
		return trace.Wrap(err)
	}
	err = agent.RequestAgentForwarding(sess)
	if err != nil {
		return trace.Wrap(err)
	}

	if err = sess.RequestSubsystem(proxySubsystemName(l.cfg.SSHUserHost, l.cfg.SSHTrustedCluster)); err != nil {
		return trace.Wrap(err)
	}
	if err := proxySession(l.context, sess); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func proxySubsystemName(userHost, cluster string) string {
	subsystem := fmt.Sprintf("proxy:%s", userHost)
	if cluster != "" {
		subsystem = fmt.Sprintf("%s@%s", subsystem, cluster)
	}
	return subsystem
}

func makeSSHClient(ctx context.Context, conn *tls.Conn, addr string, cfg *ssh.ClientConfig) (*tracessh.Client, error) {
	cc, chs, reqs, err := tracessh.NewClientConn(ctx, conn, addr, cfg)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return tracessh.NewClient(cc, chs, reqs), nil
}

func proxySession(ctx context.Context, sess *ssh.Session) error {
	stdout, err := sess.StdoutPipe()
	if err != nil {
		return trace.Wrap(err)
	}
	stdin, err := sess.StdinPipe()
	if err != nil {
		return trace.Wrap(err)
	}
	stderr, err := sess.StderrPipe()
	if err != nil {
		return trace.Wrap(err)
	}

	errC := make(chan error)
	go func() {
		defer sess.Close()
		_, err := io.Copy(os.Stdout, stdout)
		errC <- err
	}()
	go func() {
		defer sess.Close()
		_, err := io.Copy(stdin, os.Stdin)
		errC <- err
	}()
	go func() {
		defer sess.Close()
		_, err := io.Copy(os.Stderr, stderr)
		errC <- err
	}()
	var errs []error
	for i := 0; i < 3; i++ {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errC:
			if err != nil && !utils.IsOKNetworkError(err) {
				errs = append(errs, err)
			}
		}
	}
	return trace.NewAggregate(errs...)
}

// Start starts the LocalProxy.
func (l *LocalProxy) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		conn, err := l.cfg.Listener.Accept()
		if err != nil {
			if utils.IsOKNetworkError(err) {
				return nil
			}
			log.WithError(err).Errorf("Failed to accept client connection.")
			return trace.Wrap(err)
		}
		go func() {
			if err := l.handleDownstreamConnection(ctx, conn, l.cfg.SNI); err != nil {
				if utils.IsOKNetworkError(err) {
					return
				}
				log.WithError(err).Errorf("Failed to handle connection.")
			}
		}()
	}
}

// GetAddr returns the LocalProxy listener address.
func (l *LocalProxy) GetAddr() string {
	return l.cfg.Listener.Addr().String()
}

// handleDownstreamConnection proxies the downstreamConn (connection established to the local proxy) and forward the
// traffic to the upstreamConn (TLS connection to remote host).
func (l *LocalProxy) handleDownstreamConnection(ctx context.Context, downstreamConn net.Conn, serverName string) error {
	defer downstreamConn.Close()

	upstreamConn, err := tls.Dial("tcp", l.cfg.RemoteProxyAddr, &tls.Config{
		NextProtos:         l.cfg.GetProtocols(),
		InsecureSkipVerify: l.cfg.InsecureSkipVerify,
		ServerName:         serverName,
		Certificates:       l.cfg.Certs,
	})
	if err != nil {
		return trace.Wrap(err)
	}
	defer upstreamConn.Close()

	errC := make(chan error, 2)
	go func() {
		defer downstreamConn.Close()
		defer upstreamConn.Close()
		_, err := io.Copy(downstreamConn, upstreamConn)
		errC <- err
	}()
	go func() {
		defer downstreamConn.Close()
		defer upstreamConn.Close()
		_, err := io.Copy(upstreamConn, downstreamConn)
		errC <- err
	}()

	var errs []error
	for i := 0; i < 2; i++ {
		select {
		case <-ctx.Done():
			return trace.NewAggregate(append(errs, ctx.Err())...)
		case err := <-errC:
			if err != nil && !utils.IsOKNetworkError(err) {
				errs = append(errs, err)
			}
		}
	}
	return trace.NewAggregate(errs...)
}

func (l *LocalProxy) Close() error {
	l.cancel()
	if l.cfg.Listener != nil {
		if err := l.cfg.Listener.Close(); err != nil {
			return trace.Wrap(err)
		}
	}
	return nil
}

// StartAWSAccessProxy starts the local AWS CLI proxy.
func (l *LocalProxy) StartAWSAccessProxy(ctx context.Context) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			NextProtos:         l.cfg.GetProtocols(),
			InsecureSkipVerify: l.cfg.InsecureSkipVerify,
			ServerName:         l.cfg.SNI,
			Certificates:       l.cfg.Certs,
		},
	}
	proxy := &httputil.ReverseProxy{
		Director: func(outReq *http.Request) {
			outReq.URL.Scheme = "https"
			outReq.URL.Host = l.cfg.RemoteProxyAddr
		},
		Transport: tr,
	}
	err := http.Serve(l.cfg.Listener, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if err := aws.VerifyAWSSignature(req, l.cfg.AWSCredentials); err != nil {
			log.WithError(err).Errorf("AWS signature verification failed.")
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		// Requests from forward proxy have original hostnames instead of
		// localhost. Set appropriate header to keep this information.
		if addr, err := utils.ParseAddr(req.Host); err == nil && !addr.IsLocal() {
			req.Header.Set("X-Forwarded-Host", req.Host)
		}

		proxy.ServeHTTP(rw, req)
	}))
	if err != nil && !utils.IsUseOfClosedNetworkError(err) {
		return trace.Wrap(err)
	}
	return nil
}
