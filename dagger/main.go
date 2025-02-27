// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"dagger/harbor-cli/internal/dagger"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	GOLANGCILINT_VERSION = "v1.61.0"
	GO_VERSION           = "1.22.5"
	SYFT_VERSION         = "v1.9.0"
	GORELEASER_VERSION   = "v2.3.2"
)

func New(
	// Local or remote directory with source code, defaults to "./"
	// +optional
	// +defaultPath="./"
	source *dagger.Directory,
) *HarborCli {
	return &HarborCli{Source: source}
}

type HarborCli struct {
	// Local or remote directory with source code, defaults to "./"
	Source *dagger.Directory
}

// Create build of Harbor CLI for local testing and development
func (m *HarborCli) BuildDev(
	ctx context.Context,
	platform string,
) *dagger.File {
	fmt.Println("🛠️  Building Harbor-Cli with Dagger...")
	// Define the path for the binary output
	os, arch, err := parsePlatform(platform)
	if err != nil {
		log.Fatalf("Error parsing platform: %v", err)
	}
	builder := dag.Container().
		From("golang:"+GO_VERSION).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithMountedDirectory("/src", m.Source). // Ensure the source directory with go.mod is mounted
		WithWorkdir("/src").
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch)

	gitCommit, _ := builder.WithExec([]string{"git", "rev-parse", "--short", "HEAD", "--always"}).Stdout(ctx)
	buildTime := time.Now().UTC().Format(time.RFC3339)
	ldflagsArgs := fmt.Sprintf(`-X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.Version=dev
						  -X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.GoVersion=%s
						  -X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.BuildTime=%s
						  -X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.GitCommit=%s
				`, GO_VERSION, buildTime, gitCommit)
	builder = builder.WithExec([]string{
		"go", "build", "-ldflags", ldflagsArgs, "-o", "/bin/harbor-cli", "/src/cmd/harbor/main.go",
	})
	return builder.File("/bin/harbor-cli")
}

// Return list of containers for list of oses and arches
//
// FIXME: there is a bug where you cannot return a list of containers right now
// this function works as expected because it is only called by other functions but
// calling it via the CLI results in an error. That is why this into a private function for
// now so that no one calls this https://github.com/dagger/dagger/issues/8202#issuecomment-2317291483
func (m *HarborCli) build(
	ctx context.Context,
	version string,
) []*dagger.Container {
	var builds []*dagger.Container

	fmt.Println("🛠️  Building with Dagger...")
	oses := []string{"linux", "darwin", "windows"}
	arches := []string{"amd64", "arm64"}

	// temp container with git installed
	temp := dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/src", m.Source).
		// --no-cache option is to avoid caching the apk package index
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithWorkdir("/src")

	gitCommit, _ := temp.WithExec([]string{"git", "rev-parse", "--short", "HEAD", "--always"}).Stdout(ctx)
	buildTime := time.Now().UTC().Format(time.RFC3339)
	ldflagsArgs := fmt.Sprintf(`-X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.Version=%s
						  -X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.GoVersion=%s
						  -X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.BuildTime=%s
						  -X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.GitCommit=%s
				`, version, GO_VERSION, buildTime, gitCommit)

	for _, goos := range oses {
		for _, goarch := range arches {
			bin_path := fmt.Sprintf("build/%s/%s/", goos, goarch)
			builder := dag.Container().
				From("golang:"+GO_VERSION+"-alpine").
				WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
				WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
				WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
				WithEnvVariable("GOCACHE", "/go/build-cache").
				WithMountedDirectory("/src", m.Source).
				WithWorkdir("/src").
				WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch).
				WithExec([]string{"go", "build", "-ldflags", ldflagsArgs, "-o", bin_path + "harbor", "/src/cmd/harbor/main.go"}).
				WithWorkdir(bin_path).
				WithExec([]string{"ls"}).
				WithEntrypoint([]string{"./harbor"})

			builds = append(builds, builder)
		}
	}
	return builds
}

// Executes Linter and writes results to a file golangci-lint.report
func (m *HarborCli) LintReport(ctx context.Context) *dagger.File {
	report := "golangci-lint.report"
	return m.lint(ctx).WithExec([]string{
		"golangci-lint", "run", "-v",
		"--out-format", "github-actions:" + report,
		"--issues-exit-code", "0",
	}).File(report)
}

// Lint Run the linter golangci-lint
func (m *HarborCli) Lint(ctx context.Context) (string, error) {
	return m.lint(ctx).WithExec([]string{"golangci-lint", "run"}).Stderr(ctx)
}

func (m *HarborCli) lint(_ context.Context) *dagger.Container {
	fmt.Println("👀 Running linter and printing results to file golangci-lint.txt.")
	linter := dag.Container().
		From("golangci/golangci-lint:"+GOLANGCILINT_VERSION+"-alpine").
		WithMountedCache("/lint-cache", dag.CacheVolume("/lint-cache")).
		WithEnvVariable("GOLANGCI_LINT_CACHE", "/lint-cache").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src")
	return linter
}

// PublishImage publishes a container image to a registry with a specific tag and signs it using Cosign.
func (m *HarborCli) PublishImage(
	ctx context.Context,
	registry, registryUsername string,
	// +optional
	// +default=["latest"]
	imageTags []string,
	registryPassword *dagger.Secret,
) []string {
	version := getVersion(imageTags)
	builders := m.build(ctx, version)
	releaseImages := []*dagger.Container{}

	for i, tag := range imageTags {
		imageTags[i] = strings.TrimSpace(tag)
		if strings.HasPrefix(imageTags[i], "v") {
			imageTags[i] = strings.TrimPrefix(imageTags[i], "v")
		}
	}
	fmt.Printf("provided tags: %s\n", imageTags)

	for _, builder := range builders {
		os, _ := builder.EnvVariable(ctx, "GOOS")
		arch, _ := builder.EnvVariable(ctx, "GOARCH")

		if os != "linux" {
			continue
		}

		ctr := dag.Container(dagger.ContainerOpts{Platform: dagger.Platform(os + "/" + arch)}).
			From("alpine:latest").
			WithWorkdir("/").
			WithFile("/harbor", builder.File("./harbor")).
			WithExec([]string{"ls", "-al"}).
			WithExec([]string{"./harbor", "version"}).
			WithEntrypoint([]string{"/harbor"})
		releaseImages = append(releaseImages, ctr)
	}

	imageAddrs := []string{}
	for _, imageTag := range imageTags {
		addr, err := dag.Container().WithRegistryAuth(registry, registryUsername, registryPassword).
			Publish(ctx,
				fmt.Sprintf("%s/%s/harbor-cli:%s", registry, "harbor-cli", imageTag),
				dagger.ContainerPublishOpts{PlatformVariants: releaseImages},
			)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Published image address: %s\n", addr)
		imageAddrs = append(imageAddrs, addr)
	}
	return imageAddrs
}

// SnapshotRelease Create snapshot non OCI artifacts with goreleaser
func (m *HarborCli) SnapshotRelease(ctx context.Context) *dagger.Directory {
	return m.goreleaserContainer().
		WithExec([]string{"goreleaser", "release", "--snapshot", "--clean", "--skip", "validate"}).
		Directory("/src/dist")
}

// Release Create release with goreleaser
func (m *HarborCli) Release(ctx context.Context, githubToken *dagger.Secret) {
	goreleaser := m.goreleaserContainer().
		WithSecretVariable("GITHUB_TOKEN", githubToken).
		WithExec([]string{"goreleaser", "release", "--clean"})
	_, err := goreleaser.Stderr(ctx)
	if err != nil {
		log.Printf("Error occured during release: %s", err)
		return
	}
	log.Println("Release tasks completed successfully 🎉")
}

// Return a container with the goreleaser binary mounted and the source directory mounted.
func (m *HarborCli) goreleaserContainer() *dagger.Container {
	// Export the syft binary from the syft container as a file to generate SBOM
	syft := dag.Container().
		From(fmt.Sprintf("anchore/syft:%s", SYFT_VERSION)).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("syft-gomod")).
		File("/syft")

	return dag.Container().
		From(fmt.Sprintf("goreleaser/goreleaser:%s", GORELEASER_VERSION)).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithFile("/bin/syft", syft).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithEnvVariable("TINI_SUBREAPER", "true")
}

// Generate CLI Documentation and return the directory containing the generated files
func (m *HarborCli) RunDoc(ctx context.Context) *dagger.Directory {
	return dag.Container().
		From("golang:"+GO_VERSION+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src/doc").
		WithExec([]string{"go", "run", "doc.go"}).
		WithExec([]string{"go", "run", "./man-docs/man_doc.go"}).
		WithWorkdir("/src").Directory("/src/doc")
}

// Executes Go tests
func (m *HarborCli) Test(ctx context.Context) (string, error) {
	test := dag.Container().
		From("golang:"+GO_VERSION+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"go", "test", "-v", "./..."})
	return test.Stdout(ctx)
}

// Executes Go tests and returns TestReport in json file
func (m *HarborCli) TestReport(ctx context.Context) *dagger.File {
	reportName := "TestReport.json"
	test := dag.Container().
		From("golang:"+GO_VERSION+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithExec([]string{"go", "install", "gotest.tools/gotestsum@latest"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"gotestsum", "--jsonfile", reportName})

	return test.File(reportName)
}

// Parse the platform string into os and arch
func parsePlatform(platform string) (string, string, error) {
	parts := strings.Split(platform, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid platform format: %s. Should be os/arch. E.g. darwin/amd64", platform)
	}
	return parts[0], parts[1], nil
}

func getVersion(tags []string) string {
	for _, tag := range tags {
		if strings.HasPrefix(tag, "v") {
			return tag
		}
	}
	return "latest"
}

// PublishImageAndSign builds and publishes container images to a registry with a specific tags and then signs and attests them with their SBOM using Cosign.
func (m *HarborCli) PublishImageAndSign(
	ctx context.Context,
	registry string,
	registryUsername string,
	registryPassword *dagger.Secret,
	imageTags []string,
	// +optional
	githubToken *dagger.Secret,
	// +optional
	actionsIdTokenRequestToken *dagger.Secret,
	// +optional
	actionsIdTokenRequestUrl string,
) (string, error) {
	imageAddrs := m.PublishImage(ctx, registry, registryUsername, imageTags, registryPassword)

    sbom, err := m.GenerateSBOM(ctx)
    if err != nil {
        return "", fmt.Errorf("failed to generate SBOM: %w", err)
    }

	_, err = m.Sign(
		ctx,
		githubToken,
		actionsIdTokenRequestUrl,
		actionsIdTokenRequestToken,
		registryUsername,
		registryPassword,
		imageAddrs[0],
	)

	if err != nil {
		return "", fmt.Errorf("failed to sign image: %w", err)
	}

	err = m.AttestImage(
        ctx,
        registryUsername,
        registryPassword,
        imageAddrs[0],
        sbom,
    )
    if err != nil {
        return "", fmt.Errorf("failed to attest image: %w", err)
    }


	fmt.Printf("Signed image: %s\n", imageAddrs)
	return imageAddrs[0], nil
}

// GenerateSBOM generates an SBOM from the go.mod file using Syft
func (m *HarborCli) GenerateSBOM(ctx context.Context) (string, error) {
    sbom := dag.Container().
        From(fmt.Sprintf("anchore/syft:%s", SYFT_VERSION)).
        WithMountedDirectory("/src", m.Source).
        WithWorkdir("/src").
        WithExec([]string{"syft", "packages", "go.mod", "-o", "spdx-json"})
    
    return sbom.Stdout(ctx)
}

// Sign signs a container image using Cosign, works also with GitHub Actions
func (m *HarborCli) Sign(ctx context.Context,
	// +optional
	githubToken *dagger.Secret,
	// +optional
	actionsIdTokenRequestUrl string,
	// +optional
	actionsIdTokenRequestToken *dagger.Secret,
	registryUsername string,
	registryPassword *dagger.Secret,
	imageAddr string,
) (string, error) {
	registryPasswordPlain, _ := registryPassword.Plaintext(ctx)

	cosing_ctr := dag.Container().From("cgr.dev/chainguard/cosign")

	// If githubToken is provided, use it to sign the image
	if githubToken != nil {
		if actionsIdTokenRequestUrl == "" || actionsIdTokenRequestToken == nil {
			return "", fmt.Errorf("actionsIdTokenRequestUrl (exist=%s) and actionsIdTokenRequestToken (exist=%t) must be provided when githubToken is provided", actionsIdTokenRequestUrl, actionsIdTokenRequestToken != nil)
		}
		fmt.Printf("Setting the ENV Vars GITHUB_TOKEN, ACTIONS_ID_TOKEN_REQUEST_URL, ACTIONS_ID_TOKEN_REQUEST_TOKEN to sign with GitHub Token")
		cosing_ctr = cosing_ctr.WithSecretVariable("GITHUB_TOKEN", githubToken).
			WithEnvVariable("ACTIONS_ID_TOKEN_REQUEST_URL", actionsIdTokenRequestUrl).
			WithSecretVariable("ACTIONS_ID_TOKEN_REQUEST_TOKEN", actionsIdTokenRequestToken)
	}

	return cosing_ctr.WithSecretVariable("REGISTRY_PASSWORD", registryPassword).
		WithExec([]string{"cosign", "env"}).
		WithExec([]string{
			"cosign", "sign", "--yes", "--recursive",
			"--registry-username", registryUsername,
			"--registry-password", registryPasswordPlain,
			imageAddr,
			"--timeout", "1m",
		}).Stdout(ctx)
}

// AttestImage attests the image with the SBOM using Cosign
func (m *HarborCli) AttestImage(
    ctx context.Context,
    registryUsername string,
    registryPassword *dagger.Secret,
    imageAddr string,
    sbom string,
) error {
    registryPasswordPlain, _ := registryPassword.Plaintext(ctx)

    cosign_ctr := dag.Container().From("cgr.dev/chainguard/cosign")

    sbomFile := cosign_ctr.WithNewFile("/tmp/sbom.json", sbom)

    _, err := sbomFile.WithSecretVariable("REGISTRY_PASSWORD", registryPassword).
        WithExec([]string{
            "cosign", "attest", "--yes", "--type", "spdx",
            "--registry-username", registryUsername,
            "--registry-password", registryPasswordPlain,
            "--predicate", "/tmp/sbom.json",
            imageAddr,
            "--timeout", "1m",
        }).Stdout(ctx)

    return err
}
