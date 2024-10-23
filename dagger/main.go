package main

import (
	"context"
	"dagger/harbor-cli/internal/dagger"
	"fmt"
	"log"

	platformFormat "github.com/containerd/platforms"
)

const (
	GOLANGCILINT_VERSION = "v1.61.0"
	GO_VERSION           = "1.22.5"
	SYFT_VERSION         = "v1.9.0"
	GORELEASER_VERSION   = "v2.1.0"
	APP_NAME             = "dagger-harbor-cli"
	PUBLISH_ADDRESS      = "demo.goharbor.io/library/harbor-cli:0.0.3"
)

type HarborCli struct{}

func (m *HarborCli) Build(
	ctx context.Context,
	// +optional
	// +defaultPath="./"
	source *dagger.Directory) *dagger.Directory {

	fmt.Println("🛠️  Building with Dagger...")
	oses := []string{"linux", "darwin", "windows"}
	arches := []string{"amd64", "arm64"}
	outputs := dag.Directory()
	for _, goos := range oses {
		for _, goarch := range arches {
			bin_path := fmt.Sprintf("build/%s/%s/", goos, goarch)
			builder := dag.Container().
				From("golang:"+GO_VERSION+"-alpine").
				WithMountedDirectory("/src", source).
				WithWorkdir("/src").
				WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
				WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
				WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
				WithEnvVariable("GOCACHE", "/go/build-cache").
				WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch).
				WithExec([]string{"go", "build", "-o", bin_path + "harbor", "/src/cmd/harbor/main.go"})
			// Get reference to build output directory in container
			outputs = outputs.WithDirectory(bin_path, builder.Directory(bin_path))
		}
	}
	return outputs
}

// Builds the Go binary for the specified platforms, like '--platform "linux/amd64"'
func (m *HarborCli) BuildDev(
	ctx context.Context,
	// +optional
	// +defaultPath="./"
	source *dagger.Directory,
	platform string) *dagger.Directory {

	fmt.Println("🛠️  Building Go Binary for the specified platforms with Dagger...")

	os := platformFormat.MustParse(platform).OS
	arch := platformFormat.MustParse(platform).Architecture

	bin_path := fmt.Sprintf("build/%s/%s/", os, arch)
	outputs := dag.Directory()
	builder := dag.Container().
				From("golang:"+GO_VERSION+"-alpine").
				WithMountedDirectory("/src", source).
				WithWorkdir("/src").
				WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
				WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
				WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
				WithEnvVariable("GOCACHE", "/go/build-cache").
				WithEnvVariable("GOOS", os).
				WithEnvVariable("GOARCH", arch).
				WithExec([]string{"go", "build", "-o", bin_path + "harbor", "/src/cmd/harbor/main.go"})
			// Get reference to build output directory in container
			outputs = outputs.WithDirectory(bin_path, builder.Directory(bin_path))
	return outputs
}

func (m *HarborCli) Lint(
	ctx context.Context,
	// +optional
	// +defaultPath="./"
	source *dagger.Directory,
) *dagger.Container {
	fmt.Println("👀 Running linter with Dagger...")
	return dag.Container().
		From("golangci/golangci-lint:"+GOLANGCILINT_VERSION+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "--timeout", "5m"})

}

func (m *HarborCli) PullRequest(ctx context.Context,
	// +optional
	// +defaultPath="./"
	source *dagger.Directory,
	githubToken string) {
	goreleaser := goreleaserContainer(source, githubToken).WithExec([]string{"release", "--snapshot", "--clean"})
	_, err := goreleaser.Stderr(ctx)
	if err != nil {
		log.Printf("❌ Error occured during snapshot release for the recently merged pull-request: %s", err)
		return
	}
	log.Println("Pull-Request tasks completed successfully 🎉")
}

func (m *HarborCli) Release(
	ctx context.Context,
	// +optional
	// +defaultPath="./"
	source *dagger.Directory,
	githubToken string) {
	goreleaser := goreleaserContainer(source, githubToken).WithExec([]string{"release", "--clean"})
	_, err := goreleaser.Stderr(ctx)
	if err != nil {
		log.Printf("Error occured during release: %s", err)
		return
	}
	log.Println("Release tasks completed successfully 🎉")
}

func (m *HarborCli) PublishImage(
	ctx context.Context,
	// +optional
	// +defaultPath="./"
	source *dagger.Directory,
	cosignKey *dagger.Secret,
	cosignPassword string,
	regUsername string,
	regPassword string,
) string {

	builder := m.Build(ctx, source)
	// Create a minimal cli_runtime container
	cli_runtime := dag.Container().
		From("alpine:latest").
		WithWorkdir("/root/").
		WithFile("/root/harbor", builder.File("/")).
		WithEntrypoint([]string{"./harbor"})

	addr, _ := cli_runtime.Publish(ctx, PUBLISH_ADDRESS)
	cosign_password := dag.SetSecret("cosign_password", cosignPassword)
	regpassword := dag.SetSecret("reg_password", regPassword)
	_, err := dag.Cosign().Sign(ctx, cosignKey, cosign_password, []string{addr}, dagger.CosignSignOpts{RegistryUsername: regUsername, RegistryPassword: regpassword})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Published to %s 🎉\n", addr)
	return addr
}

func goreleaserContainer(directoryArg *dagger.Directory, githubToken string) *dagger.Container {
	token := dag.SetSecret("github_token", githubToken)

	// Export the syft binary from the syft container as a file to generate SBOM
	syft := dag.Container().From(fmt.Sprintf("anchore/syft:%s", SYFT_VERSION)).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("gomod")).
		File("/syft")
	return dag.Container().From(fmt.Sprintf("goreleaser/goreleaser:%s", GORELEASER_VERSION)).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("gomod")).
		WithFile("/bin/syft", syft).
		WithMountedDirectory("/src", directoryArg).WithWorkdir("/src").
		WithEnvVariable("TINI_SUBREAPER", "true").
		WithSecretVariable("GITHUB_TOKEN", token)
}
