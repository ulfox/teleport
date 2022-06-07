---
authors: Tiago Silva (tiago.silva@goteleport.com)
state: draft
---

# RFD XX - AWS Kubernetes Cluster Automatic discovery

## What

Proposes the implementation for Teleport's Kubernetes service agent to automatically discover and enroll EKS clusters.

### Related issues

- [#12048](https://github.com/gravitational/teleport/issues/12048)

## Why

Currently, when an operator wants to configure a new Kubernetes cluster in the Teleport, he can opt for these two methods:

- Helm chart: when using this method, the operator has to install `helm` binary, configure the Teleport Helm repo, and check all the configurable values (high availability, roles, apps, storage...). After that, he must create a Teleport invitation token using `tctl` and finally do the Helm install with the desired configuration.

- `Kubeconfig`: when using the `kubeconfig` procedure, the operator has to connect to the cluster with his credentials and generate a new service account for Teleport with the desired RBAC permissions. After the creation, he must extract the service account token and generate a `kubeconfig` file with the cluster CA and API server. If other clusters were already added, the operator has to merge multiple `kubeconfig` files into a single [kubeconfig][kubeconfig]. Finally, he must configure the `kubeconfig` location in Teleport config under `kubernetes_service.kubeconfig_file`.

Both processes described above are error-prone and can be tedious if the number of clusters to add to Teleport is high.

This document describes the changes required for Teleport to be able to identify the clusters based on regions and desired tags. If the clusters matched the filtering criteria they will be automatically plugged into the cluster. Once the Kubernetes is deleted, or it no longer satisfies the discovery conditions, Teleport automatically removes it from its lists.

### Scope

This RFD focuses only on AWS EKS clusters. Similar ideas will be explored in the future for GCP's GKE and Azure's AKS.

## Details

### AWS EKS discovery and IAM

AWS API has a method that allows the listing of every EKS cluster by calling [`eks:ListClusters`][listclusters] endpoint. This endpoint returns every EKS cluster name to which the user has access. The response has no other details besides the cluster name. So, for Teleport to be able to extract the cluster details such as CA and API endpoint, it has to make an extra call, per cluster, to [`eks:DescribeCluster`][descclusters]. The [`eks:DescribeCluster`][descclusters] request returns information like cluster tags used for filtering the discoveries, cluster API endpoint, and cluster Certificate Authority data.

The necessary IAM permissions required for calling those two methods are:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "eks:DescribeCluster",
                "eks:ListClusters"
            ],
            "Resource": "*" # can be limited
        }
    ]
}         
```

Besides the previous IAM policies required for the discovery service, both discovery and EKS node pool must have the `sts:GetCallerIdentity` IAM permission. The discovery requires it in order to generate an access token to Kubernetes API and the nodes for using the [Teleport's IAM connector][telepportiamrole] **Only necessary if helm is chosen**.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "sts:GetCallerIdentity"
            ],
            "Resource": "*"
        }
    ]
}
```

With cluster details such as API endpoint and CA, Teleport creates a session to gain access to the cluster. The access is given by a token generated using the [AWS IAM authenticator][awsiamauthenticator] project. This project generates a short-lived user access token by mapping the user IAM credentials/role into [Kubernetes RBAC credentials][iamroleatch] ([eksctl IAM Mappings][iameksctl]). To be able to connect to Kubernetes API it is required that the `configmap/aws-auth` maps the Teleport role into a Kubernetes group.

```yaml
apiVersion: v1
data:
  mapRoles: |
    - groups:
      - system:masters
      rolearn: arn:aws:iam::222222222222:role/teleport-role
      username: system:teleport
...
```

This means that this mapping has to exist for each cluster that the operator wants to be discovered and to do that, the operator has to change the `configmap/aws-auth`. Teleport will require access to `system:masters` (TBD the real perms required depending on the installation procedure).

The Teleport configuration for automatic AWS EKS discovery will have the following structure:

```yaml
kubernetes_service:
  enabled: yes
  aws:
  - regions: ["us-west-1"]
    tags:
      "env": "prod"
    apps: # if apps is defined, the agent registers also the APP list
      - name: grafana
        uri: http://grafana.monitoring:3000
      - name: prometheus
        uri: http://prometheus.monitoring:3000
    
    accounts:
    - aws_account: "222222222222"
      aws_roles:
      - "arn:aws:iam::222222222222:role/teleport-DescribeEKS-role"
      - "arn:aws:iam::222222222222:role/teleport-role"

  - regions: ["us-east-1", "us-east-2"]
    tags:
      "env": "stage"

```

### Cluster configuration

#### Possibility 1: Helm chart

Teleport has a join method, IAM join, that allows Teleport agents and Proxies to join a Teleport cluster without sharing any secrets when they are running in AWS.

The IAM join method is available to any Teleport agent running anywhere with access to IAM credentials, such as an EC2 instance that is part of an EKS cluster. This method allows any resource that fulfills the defined criteria to be able to join automatically into the Teleport cluster.

Every agent requires access to `sts:GetCallerIdentity` in order to use the IAM method. No other IAM policy or permissions are needed.

To avoid attempting to reinstall the Teleport agent in clusters where it is already present, the Teleport config will include a static label indicating it was created via auto-discovery.

```json
{
    "kind": "kube_service",
    "version": "v2",
    "metadata": {
        "name": "{uuid}",
        # the same token used for AWS EC2 instances auto-discovery
        "teleport.dev/discovered-node": "yes"
    },
    "spec": {
        "addr": "..",
        "hostname": "",
        "version": "{VERSION}",
        "kube_clusters": [
            {
                "name": "{CLUSTER_NAME}",
                "static_labels": {
                   {CLUSTER_LABELS[:]}
                }
            }
        ]
    }
}
```

To configure the IAM joining token method, the operator has to define the IAM token spec.

```yaml

kind: token
version: v2
metadata:
  # the token name is not a secret because instances must prove that they are
  # running in your AWS account to use this token
  name: kube-iam-token
  # set a long expiry time, the default for tokens is only 30 minutes
  expires: "3000-01-01T00:00:00Z"
spec:
  # use the minimal set of roles required
  roles: [Kube,[App]]

  # set the join method allowed for this token
  join_method: iam

  allow:
  # specify the AWS account which nodes may join from
  - aws_account: "111111111111"
  # multiple allow rules are supported
  - aws_account: "222222222222"
  # aws_arn is optional and allows you to restrict the IAM role of joining nodes
  - aws_account: "333333333333"
    aws_arn: "arn:aws:sts::333333333333:assumed-role/teleport-node-role/i-*"
```

Once Teleport has discovered a cluster and granted access to its API, it installs the Helm Agent chart via [Helm library][helmlib]. The Teleport Helm chart has to be updated to support IAM joining token.
Teleport discovery has to define the correct values for the Helm chart installation and execute them into the cluster. After that, the deployment is done.

##### Limitations

- If helm chart values are updated the discovery code must also be updated accordingly
- Requires Kube Secrets backend storage.

#### Possibility 2: Kubeconfig

Teleport supports loading clusters from Kubernetes Kubeconfig during the startup procedure. The Kubeconfig supports multiple Kubernetes clusters and has 3 main sections, `clusters`, `contexts` and `users`:

```yaml

apiVersion: v1
clusters:
- cluster:
    certificate-authority: ca.crt
    server: https://127.0.0.1:52181
  name: minikube
contexts:
- context:
    cluster: minikube
    namespace: default
    user: minikube
  name: minikube
kind: Config
users:
- name: minikube
  user:
    client-certificate: client.crt
    client-key: client.key
```

- `clusters`: defines the `cluster` API and certificate authority.
- `users`: describes `user` credentials.
- `contexts`: links a `cluster` and `user` definition.

Teleport discovery runs the cluster discovery mechanism every 60 seconds and will check if new clusters appeared or if any were deleted.

To enroll a new cluster, Teleport discovery must create a new service account with the following `ClusterRole` associated.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: teleport-role
rules:
- apiGroups:
  - ""
  resources:
  - users
  - groups
  - serviceaccounts
  verbs:
  - impersonate
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - get
- apiGroups:
  - "authorization.k8s.io"
  resources:
  - selfsubjectaccessreviews
  - selfsubjectrulesreviews
  verbs:
  - create
```

Once the `ClusterRoleBinding` is defined, the SA token has to be extracted to build the `user` section for kubeconfig.

Once all clusters are enrolled, and Teleport has the credentials for all of them, Teleport reloads the Kubeconfig in order to update its state.

##### Limitations

- Teleport has to check not only for new EKS clusters but also for the ones deletes in order to remove them.
- More code has to be added into Teleport in order to support dynamic credentials extraction from kubeconfig.


# Links
[kubeconfig]: https://goteleport.com/docs/kubernetes-access/guides/standalone-teleport/#step-12-generate-a-kubeconfig
[listclusters]: https://docs.aws.amazon.com/eks/latest/APIReference/API_ListClusters.html
[descclusters]: https://docs.aws.amazon.com/eks/latest/APIReference/API_DescribeCluster.html
[awsiamauthenticator]: https://github.com/kubernetes-sigs/aws-iam-authenticator
[telepportiamrole]: https://goteleport.com/docs/setup/guides/joining-nodes-aws-iam/
[iamroleatch]: https://docs.aws.amazon.com/eks/latest/userguide/add-user-role.html
[iameksctl]: https://eksctl.io/usage/iam-identity-mappings/
[helmlib]: https://github.com/helm/helm/tree/main/pkg
1. https://goteleport.com/docs/kubernetes-access/guides/standalone-teleport/#step-12-generate-a-kubeconfig
2. https://docs.aws.amazon.com/eks/latest/APIReference/API_ListClusters.html
3. https://docs.aws.amazon.com/eks/latest/APIReference/API_DescribeCluster.html
4. https://github.com/kubernetes-sigs/aws-iam-authenticator
5. https://goteleport.com/docs/setup/guides/joining-nodes-aws-iam/
6. https://docs.aws.amazon.com/eks/latest/userguide/add-user-role.html
7. https://eksctl.io/usage/iam-identity-mappings/
8. https://github.com/helm/helm/tree/main/pkg