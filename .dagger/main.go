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
	GOLANGCILINT_VERSION = "v2.1.2"
	GO_VERSION           = "1.24.4"
	GORELEASER_VERSION   = "v2.8.2"
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

// Executes Linter and writes results to a file golangci-lint.report
func (m *HarborCli) LintReport(ctx context.Context) *dagger.File {
	report := "golangci-lint.report"
	return m.lint(ctx).WithExec([]string{
		"golangci-lint", "run", "-v",
		"--output.tab.path=" + report,
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

// PublishImage publishes a container image to a registry with a specific tags using goreleaser-generated binaries
func (m *HarborCli) PublishImage(
	ctx context.Context,
	registry string,
	registryUsername string,
	// +optional
	// +default=["latest"]
	imageTags []string,
	registryPassword *dagger.Secret,
	// +optional
	githubToken *dagger.Secret,
	// +optional
	snapshot bool,
) []string {
	version := getVersion(imageTags)
	for i, tag := range imageTags {
		imageTags[i] = strings.TrimSpace(tag)
		if strings.HasPrefix(imageTags[i], "v") {
			imageTags[i] = strings.TrimPrefix(imageTags[i], "v")
		}
	}
	fmt.Printf("🚀 Publishing Harbor-Cli image to %s with processed tags: %s...\n", registry, imageTags)

	// Use goreleaser to generate binaries for multiple platforms
	distDir := m.RunGoreleaser(ctx, githubToken, snapshot)

	// Get current time for image creation timestamp
	creationTime := time.Now().UTC().Format(time.RFC3339)

	platforms := []struct {
		goos, goarch, distPath string
		platform               string
	}{
		{"linux", "amd64", "harbor_linux_amd64_v1", "linux/amd64"},
		{"linux", "arm64", "harbor_linux_arm64_v8.0", "linux/arm64"},
	}

	// Create containers for each platform
	var releaseImages []*dagger.Container
	for _, p := range platforms {
		binaries := distDir.Directory(p.distPath)

		ctr := dag.Container(dagger.ContainerOpts{Platform: dagger.Platform(p.platform)}).
			From("alpine:latest").
			WithDirectory("/tmp/dist", binaries).
			WithExec([]string{"sh", "-c", "cp /tmp/dist/harbor /usr/local/bin/harbor && chmod +x /usr/local/bin/harbor"}).
			WithExec([]string{"sh", "-c", "mkdir -p /usr/share/doc/harbor && find /tmp/dist -name '*.sbom' -exec cp {} /usr/share/doc/harbor/ \\; || true"}).
			WithEntrypoint([]string{"/usr/local/bin/harbor"}).
			WithExec([]string{"ls", "-al", "/usr/local/bin/"}).
			WithExec([]string{"/usr/local/bin/harbor", "version"}).
			WithLabel("org.opencontainers.image.created", creationTime).
			WithLabel("org.opencontainers.image.description", "Harbor CLI - A command-line interface for CNCF Harbor, the cloud native registry!").
			WithLabel("io.artifacthub.package.readme-url", "https://raw.githubusercontent.com/goharbor/harbor-cli/main/README.md").
			WithLabel("org.opencontainers.image.source", "https://github.com/goharbor/harbor-cli").
			WithLabel("org.opencontainers.image.version", version).
			WithLabel("io.artifacthub.package.license", "Apache-2.0").
			WithLabel("org.opencontainers.image.documentation", "SBOM available at /usr/share/doc/harbor/")

		releaseImages = append(releaseImages, ctr)
	}

	var imageAddrs []string
	for _, tag := range imageTags {
		imageAddr := fmt.Sprintf("%s/%s/harbor-cli:%s", registry, registryUsername, tag)

		publishedImage, err := dag.Container().
			WithRegistryAuth(registry, registryUsername, registryPassword).
			Publish(ctx, imageAddr, dagger.ContainerPublishOpts{
				PlatformVariants: releaseImages,
			})

		if err != nil {
			log.Printf("Error publishing image %s: %v", imageAddr, err)
			continue
		}

		fmt.Printf("Published image address: %s\n", publishedImage)
		imageAddrs = append(imageAddrs, publishedImage)
	}

	return imageAddrs
}

// RunGoreleaser runs goreleaser to generate binaries
func (m *HarborCli) RunGoreleaser(
	ctx context.Context,
	// +optional
	githubToken *dagger.Secret,
	// +optional
	snapshot bool,
) *dagger.Directory {
	fmt.Println("🚀 Running goreleaser to generate binaries...")

	goreleaser := dag.Container().
		From("goreleaser/goreleaser:"+GORELEASER_VERSION).
		WithMountedDirectory("/src", m.Source).
		WithMountedDirectory("/src/.git", m.Source.Directory(".git")).
		WithWorkdir("/src")

	args := []string{"goreleaser", "release", "--clean"}

	if snapshot {
		args = append(args, "--snapshot")
	}

	if githubToken != nil {
		goreleaser = goreleaser.WithSecretVariable("GITHUB_TOKEN", githubToken)
	}

	goreleaser = goreleaser.WithExec(args)

	return goreleaser.Directory("/src/dist")
}

// PublishImageAndSign builds and publishes container images to a registry with specific tags and signs them using Cosign.
func (m *HarborCli) PublishImageAndSign(
	ctx context.Context,
	registry string,
	registryUsername string,
	// +optional
	// +default=["latest"]
	imageTags []string,
	registryPassword *dagger.Secret,
	// +optional
	githubToken *dagger.Secret,
	// +optional
	actionsIdTokenRequestToken *dagger.Secret,
	// +optional
	actionsIdTokenRequestUrl string,
	// +optional
	snapshot bool,
) (string, error) {
	fmt.Println("🚀 Starting PublishImageAndSign...")

	// First publish the image
	fmt.Println("📦 Publishing image...")
	imageAddrs := m.PublishImage(
		ctx,
		registry,
		registryUsername,
		imageTags,
		registryPassword,
		githubToken,
		snapshot,
	)

	if len(imageAddrs) == 0 {
		return "", fmt.Errorf("no images were published")
	}
	fmt.Printf("✅ Published image: %s\n", imageAddrs[0])

	// If no registry password is provided, skip signing
	if registryPassword == nil {
		fmt.Println("⚠️  No registry password provided, skipping image signing")
		return imageAddrs[0], nil
	}

	// Then sign the first image
	fmt.Println("🔏 Starting image signing process...")
	signedImage, err := m.Sign(
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

	fmt.Printf("✅ Successfully signed image: %s\n", signedImage)
	return signedImage, nil
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

// SnapshotRelease Create snapshot non OCI artifacts with goreleaser
func (m *HarborCli) SnapshotRelease(ctx context.Context) *dagger.Directory {
	return m.RunGoreleaser(ctx, nil, true)
}

// Release Create release with goreleaser
func (m *HarborCli) Release(ctx context.Context, githubToken *dagger.Secret) {
	distDir := m.RunGoreleaser(ctx, githubToken, false)
	_, err := distDir.Entries(ctx)
	if err != nil {
		log.Printf("Error occurred during release: %s", err)
		return
	}
	if len(error) > 0 {
		log.Printf("Error occured while release: %s", err)
		return
	}
	log.Println("Release tasks completed successfully 🎉")
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
// TestReport executes Go tests and returns only the JSON report file
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
		WithExec([]string{"gotestsum", "--jsonfile", reportName, "./..."})

	return test.File(reportName)
}

func (m *HarborCli) TestCoverage(ctx context.Context) *dagger.File {
	coverage := "coverage.out"
	test := dag.Container().
		From("golang:"+GO_VERSION+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithExec([]string{"go", "install", "gotest.tools/gotestsum@latest"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"gotestsum", "--", "-coverprofile=" + coverage, "./..."})

	return test.File(coverage)
}

// TestCoverageReport processes coverage data and returns a formatted markdown report
func (m *HarborCli) TestCoverageReport(ctx context.Context) *dagger.File {
	coverageFile := "coverage.out"
	reportFile := "coverage-report.md"
	test := dag.Container().
		From("golang:"+GO_VERSION+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "bc"}).
		WithExec([]string{"go", "test", "-coverprofile=" + coverageFile, "./..."})
	return test.WithExec([]string{"sh", "-c", `
        echo "<h2> 📊 Test Coverage Results</h2>" > ` + reportFile + `
        if [ ! -f "` + coverageFile + `" ]; then
            echo "<p>❌ Coverage file not found!</p>" >> ` + reportFile + `
            exit 1
        fi
        total_coverage=$(go tool cover -func=` + coverageFile + ` | grep total: | grep -Eo '[0-9]+\.[0-9]+')
        echo "DEBUG: Total coverage is $total_coverage" >&2
        if (( $(echo "$total_coverage >= 80.0" | bc -l) )); then
            emoji="✅"
        elif (( $(echo "$total_coverage >= 60.0" | bc -l) )); then
            emoji="⚠️"
        else
            emoji="❌"
        fi
		echo "<p><b>Total coverage: $emoji $total_coverage% (Target: 80%)</b></p>" >> ` + reportFile + `
		echo "<details><summary>Detailed package coverage</summary><pre>" >> ` + reportFile + `
        go tool cover -func=` + coverageFile + ` >> ` + reportFile + `
        echo "</pre></details>" >> ` + reportFile + `
        cat ` + reportFile + ` >&2
    `}).File(reportFile)
}

// Checks for vulnerabilities using govulncheck
func (m *HarborCli) vulnerabilityCheck(ctx context.Context) *dagger.Container {
	return dag.Container().
		From("golang:"+GO_VERSION+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithExec([]string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@latest"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src")
}

// Runs a vulnerability check using govulncheck
func (m *HarborCli) VulnerabilityCheck(ctx context.Context) (string, error) {
	return m.vulnerabilityCheck(ctx).
		WithExec([]string{"govulncheck", "-show", "verbose", "./..."}).
		Stderr(ctx)
}

// Runs a vulnerability check using govulncheck and writes results to vulnerability-check.report
func (m *HarborCli) VulnerabilityCheckReport(ctx context.Context) *dagger.File {
	report := "vulnerability-check.report"
	return m.vulnerabilityCheck(ctx).
		WithExec([]string{
			"sh", "-c", fmt.Sprintf("govulncheck ./... > %s", report),
		}).File(report)
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
