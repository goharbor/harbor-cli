package main

import (
	"context"
	"dagger/harbor-cli/internal/dagger"
	"fmt"
	"log"
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

// Return list of containers for list of oses and arches
//
// FIXME: there is a bug where you cannot return a list of containers right now
// this function works as expected because it is only called by other functions but
// calling it via the CLI results in an error. That is why this into a private function for
// now so that no one calls this https://github.com/dagger/dagger/issues/8202#issuecomment-2317291483
func (m *HarborCli) build(
	ctx context.Context,
) []*dagger.Container {
	var builds []*dagger.Container

	fmt.Println("🛠️  Building with Dagger...")
	oses := []string{"linux", "darwin", "windows"}
	arches := []string{"amd64", "arm64"}
	for _, goos := range oses {
		for _, goarch := range arches {
			bin_path := fmt.Sprintf("build/%s/%s/", goos, goarch)
			builder := dag.Container().
				From("golang:"+GO_VERSION+"-alpine").
				WithMountedDirectory("/src", m.Source).
				WithWorkdir("/src").
				WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
				WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
				WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
				WithEnvVariable("GOCACHE", "/go/build-cache").
				WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch).
				WithExec([]string{"go", "build", "-o", bin_path + "harbor", "/src/cmd/harbor/main.go"}).
				WithWorkdir(bin_path).
				WithExec([]string{"ls"}).
				WithEntrypoint([]string{"./harbor"})

			builds = append(builds, builder)
		}
	}
	return builds
}

func (m *HarborCli) Lint(
	ctx context.Context,
) *dagger.Container {
	fmt.Println("👀 Running linter with Dagger...")
	return dag.Container().
		From("golangci/golangci-lint:"+GOLANGCILINT_VERSION+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "--timeout", "5m"})
}

func (m *HarborCli) PullRequest(
	ctx context.Context,
	// Github API token
	githubToken *dagger.Secret,
) {
	_, err := m.
		goreleaserContainer(githubToken).
		WithExec([]string{"release", "--snapshot", "--clean"}).
		Stderr(ctx)
	if err != nil {
		log.Printf("❌ Error occured during snapshot release for the recently merged pull-request: %s", err)
		return
	}
	log.Println("Pull-Request tasks completed successfully 🎉")
}

// Create release with goreleaser
func (m *HarborCli) Release(
	ctx context.Context,
	// Github API token
	githubToken *dagger.Secret,
) {
	goreleaser := m.goreleaserContainer(githubToken).
		WithExec([]string{"ls", "-la"}).
		WithExec([]string{"goreleaser", "release", "--clean"})

	_, err := goreleaser.Stderr(ctx)
	if err != nil {
		log.Printf("Error occured during release: %s", err)
		return
	}
	log.Println("Release tasks completed successfully 🎉")
}

// PublishImage publishes a Docker image to a registry with a specific tag and signs it using Cosign.
func (m *HarborCli) PublishImage(
	ctx context.Context,
	// the secret used for signing the image
	cosignKey *dagger.Secret,
	// the password for the cosign secret
	cosignPassword *dagger.Secret,
	// the username for the registry
	regUsername string,
	// the password for the registry
	regPassword *dagger.Secret,
	// the address of the registry to publish the image
	regAddress string,
	// the address of the registry to publish the image
	publishAddress string,
	// the version tag for the image
	tag string,
) string {
	var container *dagger.Container
	var filteredBuilders []*dagger.Container

	builders := m.build(ctx)
	if len(builders) > 0 {
		fmt.Println(len(builders))
		container = builders[0]
		builders = builders[3:6]
	}
	dir := dag.Directory()
	dir = dir.WithDirectory(".", container.Directory("."))

	// Create a minimal cli_runtime container
	cli_runtime := dag.Container().
		From("alpine:latest").
		WithWorkdir("/root/").
		WithFile("/root/harbor", dir.File("./harbor")).
		WithExec([]string{"ls"}).
		WithExec([]string{"./harbor", "--help"}).
		WithEntrypoint([]string{"./harbor"})

	for _, builder := range builders {
		if !(buildPlatform(ctx, builder) == "linux/amd64") {
			filteredBuilders = append(filteredBuilders, builder)
		}
	}

	publisher := cli_runtime.WithRegistryAuth(regAddress, regUsername, regPassword)
	// Push the versioned tag
	versionedAddress := fmt.Sprintf("%s:%s", publishAddress, tag)
	addr, err := publisher.Publish(ctx, versionedAddress, dagger.ContainerPublishOpts{PlatformVariants: filteredBuilders})
	if err != nil {
		panic(err)
	}
	// Push the latest tag
	latestAddress := fmt.Sprintf("%s:latest", publishAddress)
	addr, err = publisher.Publish(ctx, latestAddress)
	if err != nil {
		panic(err)
	}

	_, err = dag.Cosign().Sign(ctx, cosignKey, cosignPassword, []string{addr}, dagger.CosignSignOpts{RegistryUsername: regUsername, RegistryPassword: regPassword})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully published image to %s 🎉\n", addr)

	return addr
}

// Return the platform of the container
func buildPlatform(
	ctx context.Context,
	// container to check
	container *dagger.Container,
) string {
	platform, err := container.Platform(ctx)
	if err != nil {
		log.Fatalf("error getting platform", err)
	}
	return string(platform)
}

// Return a container with the goreleaser binary mounted and the source directory mounted.
func (m *HarborCli) goreleaserContainer(
	// Github API token
	githubToken *dagger.Secret,
) *dagger.Container {
	// Export the syft binary from the syft container as a file to generate SBOM
	syft := dag.Container().
		From(fmt.Sprintf("anchore/syft:%s", SYFT_VERSION)).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("syft-gomod")).
		File("/syft")

	return dag.Container().
		From(fmt.Sprintf("goreleaser/goreleaser:%s", GORELEASER_VERSION)).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("goreleaser-gomod")).
		WithFile("/bin/syft", syft).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithEnvVariable("TINI_SUBREAPER", "true").
		WithSecretVariable("GITHUB_TOKEN", githubToken)
}

// Generate CLI Documentation with doc.go and return the directory containing the generated files
func (m *HarborCli) RunDoc(
	ctx context.Context,
) *dagger.Directory {
	fmt.Println("Running doc.go file using Dagger...")
	return dag.Container().
		From("golang:latest").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod")).
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build")).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src/doc").
		WithExec([]string{"go", "run", "doc.go"}).
		WithWorkdir("/src").Directory("/src/doc")
}
