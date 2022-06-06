// Copyright 2021 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"path"
	"strings"
)

func promoteBuildPipeline() pipeline {
	aptPipeline := promoteAptPipeline()
	return aptPipeline
}

// Used for one-off migrations of older versions.
// Use cases include:
//  * We want to support another OS while providing backwards compatibility
//  * We want to support another OS version while providing backwards compatibility
//  * A customer wants to be able to install an older version via APT/YUM even if we
//      no longer support it
//  * RPM migrations after new YUM pipeline is done
func artifactMigrationPipeline() pipeline {
	migrationVersions := []string{
		// These versions were migrated as a part of the new `promoteAptPipeline`
		// "v6.2.31",
		// "v7.3.17",
		// "v7.3.18",
		// "v7.3.19",
		// "v7.3.20",
		// "v7.3.21",
		// "v8.3.3",
		// "v8.3.4",
		// "v8.3.5",
		// "v8.3.6",
		// "v8.3.7",
		// "v8.3.8",
		// "v8.3.9",
		// "v8.3.10",
		// "v8.3.11",
		// "v8.3.12",
		// "v9.0.0",
		// "v9.0.1",
		// "v9.0.2",
		// "v9.0.3",
		// "v9.0.4",
		// "v9.1.0",
		// "v9.1.1",
		// "v9.1.2",
		// "v9.1.3",
		// "v9.2.0",
		// "v9.2.1",
		// "v9.2.2",
		// "v9.2.3",
		// "v9.2.4",
		// "v9.3.0",
		// "v9.3.2",
	}
	// Pushing to this branch will trigger the listed versions to be migrated. Typically this should be
	// the branch that these changes are being committed to.
	migrationBranch := "" // "rfd/0058-package-distribution"

	aptPipeline := migrateAptPipeline(migrationBranch, migrationVersions)
	return aptPipeline
}

// This function calls the build-apt-repos tool which handles the APT portion of RFD 0058.
func promoteAptPipeline() pipeline {
	aptVolumeName := "aptrepo"
	commitName := "${DRONE_TAG}"

	p := buildBaseAptPipeline("publish-apt-new-repos", aptVolumeName, commitName)
	p.Trigger = triggerPromote
	p.Trigger.Repo.Include = []string{"gravitational/teleport"}

	steps := []step{
		{
			Name:  "Verify build is tagged",
			Image: "alpine:latest",
			Commands: []string{
				"[ -n ${DRONE_TAG} ] || (echo 'DRONE_TAG is not set. Is the commit tagged?' && exit 1)",
			},
		},
	}
	steps = append(steps, p.Steps...)
	steps = append(steps, getDroneTagVersionSteps(aptVolumeName)...)
	p.Steps = steps

	return p
}

func migrateAptPipeline(triggerBranch string, migrationVersions []string) pipeline {
	aptVolumeName := "aptrepo"
	pipelineName := "migrate-apt-new-repos"
	// DRONE_TAG is not available outside of promotion pipelines and will cause drone to fail with a
	// "migrate-apt-new-repos: bad substitution" error if used here
	commitName := "${DRONE_COMMIT}"

	// If migrations are not configured then don't run
	if triggerBranch == "" || len(migrationVersions) == 0 {
		return buildNeverTriggerPipeline(pipelineName)
	}

	p := buildBaseAptPipeline(pipelineName, aptVolumeName, commitName)
	p.Trigger = trigger{
		Repo:   triggerRef{Include: []string{"gravitational/teleport"}},
		Event:  triggerRef{Include: []string{"push"}},
		Branch: triggerRef{Include: []string{triggerBranch}},
	}

	for _, migrationVersion := range migrationVersions {
		p.Steps = append(p.Steps, getVersionSteps(migrationVersion, aptVolumeName)...)
	}

	return p
}

// Builds a pipeline that is syntactically correct but should never trigger to create
// a placeholder pipeline
func buildNeverTriggerPipeline(pipelineName string) pipeline {
	p := newKubePipeline(pipelineName)
	p.Trigger = trigger{
		Event:  triggerRef{Include: []string{"custom"}},
		Repo:   triggerRef{Include: []string{"non-existent-repository"}},
		Branch: triggerRef{Include: []string{"non-existent-branch"}},
	}

	p.Steps = []step{
		{
			Name:  "Placeholder",
			Image: "alpine:latest",
			Commands: []string{
				"echo \"This command, step, and pipeline never runs\"",
			},
		},
	}

	return p
}

// Functions that use this method should add at least:
// * a Trigger
// * Steps for checkout
func buildBaseAptPipeline(pipelineName, aptVolumeName, commit string) pipeline {
	p := newKubePipeline(pipelineName)
	p.Workspace = workspace{Path: "/go"}
	p.Volumes = []volume{
		{
			Name: aptVolumeName,
			Claim: &volumeClaim{
				Name: "drone-s3-aptrepo-pvc",
			},
		},
		volumeTmpfs,
	}
	p.Steps = []step{
		{
			Name:     "Check out code",
			Image:    "alpine/git:latest",
			Commands: aptToolCheckoutCommands(commit),
		},
	}

	return p
}

func getDroneTagVersionSteps(aptVolumeName string) []step {
	return getVersionSteps("${DRONE_TAG}", aptVolumeName)
}

// Version should start with a 'v', i.e. v1.2.3 or v9.0.1, or should be an environment var
// i.e. ${DRONE_TAG}
func getVersionSteps(version, aptVolumeName string) []step {
	artifactPath := "/go/artifacts"
	pvcMountPoint := "/mnt"

	var bucketFolder string
	switch version[0:1] {
	// If environment var
	case "$":
		// Remove the 'v' at runtime as the value isn't known at compile time
		// This will change "${SOME_VAR}" to "${SOME_VAR##v}". `version` isn't actually
		// an environment variable - it's a Drone substitution variable. See
		// https://docs.drone.io/pipeline/environment/substitution/ for details.
		bucketFolder = fmt.Sprintf("%s##v}", version[:len(version)-1])
	// If static string
	case "v":
		// Remove the 'v' at compile time as the value is known then
		bucketFolder = version[1:]
	}

	return []step{
		{
			Name:  fmt.Sprintf("Download artifacts for %q", version),
			Image: "amazon/aws-cli",
			Environment: map[string]value{
				"AWS_S3_BUCKET": {
					fromSecret: "AWS_S3_BUCKET",
				},
				"AWS_ACCESS_KEY_ID": {
					fromSecret: "AWS_ACCESS_KEY_ID",
				},
				"AWS_SECRET_ACCESS_KEY": {
					fromSecret: "AWS_SECRET_ACCESS_KEY",
				},
				"ARTIFACT_PATH": {
					raw: artifactPath,
				},
			},
			Commands: []string{
				"mkdir -pv \"$ARTIFACT_PATH\"",
				strings.Join(
					[]string{
						"aws s3 sync",
						"--no-progress",
						"--delete",
						"--exclude \"*\"",
						"--include \"*.deb*\"",
						fmt.Sprintf("s3://$AWS_S3_BUCKET/teleport/tag/%s/", bucketFolder),
						"\"$ARTIFACT_PATH\"",
					},
					" ",
				),
			},
		},
		{
			Name: fmt.Sprintf("Publish debs to APT repos for %q", version),
			// TODO set this if drongen `step` supports https://docs.drone.io/pipeline/ssh/syntax/parallelism/ in the future
			// DependsOn: []string {
			// 	"Check out code",
			// 	"Download artifacts",
			// },
			Image: "golang:1.18.1-bullseye",
			Environment: map[string]value{
				"APT_S3_BUCKET": {
					fromSecret: "APT_REPO_NEW_AWS_S3_BUCKET",
				},
				"BUCKET_CACHE_PATH": {
					// If we need to cache the bucket on the PVC for some reason in the future
					// uncomment this line
					// raw: path.Join(pvcMountPoint, "bucket-cache"),
					raw: "/tmp/bucket",
				},
				"AWS_REGION": {
					raw: "us-west-2",
				},
				"AWS_ACCESS_KEY_ID": {
					fromSecret: "APT_REPO_NEW_AWS_ACCESS_KEY_ID",
				},
				"AWS_SECRET_ACCESS_KEY": {
					fromSecret: "APT_REPO_NEW_AWS_SECRET_ACCESS_KEY",
				},
				"ARTIFACT_PATH": {
					raw: artifactPath,
				},
				"APTLY_ROOT_DIR": {
					raw: path.Join(pvcMountPoint, "aptly"),
				},
				"GNUPGHOME": {
					raw: "/tmpfs/gnupg",
				},
				"GPG_RPM_SIGNING_ARCHIVE": {
					fromSecret: "GPG_RPM_SIGNING_ARCHIVE",
				},
			},
			Commands: []string{
				"mkdir -pv -m0700 $GNUPGHOME",
				"echo \"$GPG_RPM_SIGNING_ARCHIVE\" | base64 -d | tar -xzf - -C $GNUPGHOME",
				"chown -R root:root $GNUPGHOME",
				"apt update",
				"apt install aptly tree -y",
				"cd /go/src/github.com/gravitational/teleport/build.assets/tooling",
				fmt.Sprintf("export VERSION=%q", version),
				"export RELEASE_CHANNEL=\"stable\"", // The tool supports several release channels but I'm not sure where this should be configured
				// "rm -rf \"$APTLY_ROOT_DIR\"",		// Uncomment this to completely dump the Aptly database and force a rebuild
				strings.Join(
					[]string{
						// This just makes the (long) command a little more readable
						"go run ./cmd/build-apt-repos",
						"-bucket \"$APT_S3_BUCKET\"",
						"-local-bucket-path \"$BUCKET_CACHE_PATH\"",
						"-artifact-version \"$VERSION\"",
						"-release-channel \"$RELEASE_CHANNEL\"",
						"-aptly-root-dir \"$APTLY_ROOT_DIR\"",
						"-artifact-path \"$ARTIFACT_PATH\"",
						"-log-level 4", // Set this to 5 for debug logging
					},
					" ",
				),
				"rm -rf \"$BUCKET_CACHE_PATH\"",
				"df -h \"$APTLY_ROOT_DIR\"",
			},
			Volumes: []volumeRef{
				{
					Name: aptVolumeName,
					Path: pvcMountPoint,
				},
				volumeRefTmpfs,
			},
		},
	}
}

// Note that tags are also valid here as a tag refers to a specific commit
func aptToolCheckoutCommands(commit string) []string {
	commands := []string{
		`mkdir -p /go/src/github.com/gravitational/teleport`,
		`cd /go/src/github.com/gravitational/teleport`,
		`git clone https://github.com/gravitational/${DRONE_REPO_NAME}.git .`,
		fmt.Sprintf("git checkout %q", commit),
	}
	return commands
}

func updateDocsPipeline() pipeline {
	// TODO: migrate
	return pipeline{}
}
