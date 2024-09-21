package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/goharbor/harbor-cli/internal/dagger"
)

const (
	GO_VERSION         = "1.22.5"
	SYFT_VERSION       = "v1.9.0"
	GORELEASER_VERSION = "v2.1.0"
	APP_NAME           = "dagger-harbor-cli"
	PUBLISH_ADDRESS    = "demo.goharbor.io/library/harbor-cli:0.0.3"
)

type HarborCli struct{}

func (m *HarborCli) Echo(stringArg string) string {
	return stringArg
}

// Returns a container that echoes whatever string argument is provided
func (m *HarborCli) ContainerEcho(stringArg string) *dagger.Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})

}

// Returns lines that match a pattern in the files of the provided Directory
func (m *HarborCli) GrepDir(ctx context.Context, directoryArg *dagger.Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)

}

func (m *HarborCli) LintCode(ctx context.Context, directoryArg *dagger.Directory) *dagger.Container {
	fmt.Println("👀 Running linter with Dagger...")
	return dag.Container().
		From("golangci/golangci-lint:v1.59.1-alpine").
		WithMountedDirectory("/src", directoryArg).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "--timeout", "5m"})

}

func (m *HarborCli) BuildHarbor(ctx context.Context, directoryArg *dagger.Directory) *dagger.Directory {
	fmt.Println("🛠️  Building with Dagger...")
	oses := []string{"linux", "darwin", "windows"}
	arches := []string{"amd64", "arm64"}
	outputs := dag.Directory()
	golangcont, main_go_path := fetchMainGoPath(ctx, directoryArg)

	for _, goos := range oses {
		for _, goarch := range arches {
			path := fmt.Sprintf("build/%s/%s/", goos, goarch)
			build := golangcont.WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch).
				WithExec([]string{"go", "build", "-o", path + "harbor", main_go_path})

			// Get reference to build output directory in container
			outputs = outputs.WithDirectory(path, build.Directory(path))
		}
	}
	return outputs
}

func (m *HarborCli) PullRequest(ctx context.Context, directoryArg *dagger.Directory, githubToken string) {
	goreleaser := goreleaserContainer(directoryArg, githubToken).WithExec([]string{"release", "--snapshot", "--clean"})
	_, err := goreleaser.Stderr(ctx)
	if err != nil {
		log.Printf("❌ Error occured during snapshot release for the recently merged pull-request: %s", err)
		return
	}
	log.Println("Pull-Request tasks completed successfully 🎉")
}

func (m *HarborCli) Release(ctx context.Context, directoryArg *dagger.Directory, githubToken string) {
	goreleaser := goreleaserContainer(directoryArg, githubToken).WithExec([]string{"release", "--clean"})
	_, err := goreleaser.Stderr(ctx)
	if err != nil {
		log.Printf("Error occured during release: %s", err)
		return
	}
	log.Println("Release tasks completed successfully 🎉")
}

func (m *HarborCli) DockerPublish(ctx context.Context, directoryArg *dagger.Directory, cosignKey *dagger.Secret, cosignPassword string, regUsername string, regPassword string) string {

	builder, main_go_path := fetchMainGoPath(ctx, directoryArg)
	builder = builder.WithWorkdir("/src").WithExec([]string{"go", "build", "-o", "harbor", main_go_path})

	// Create a minimal runtime container
	runtime := dag.Container().
		From("alpine:latest").
		WithWorkdir("/root/").
		WithFile("/root/harbor", builder.File("/src/harbor")).
		WithEntrypoint([]string{"./harbor"})

	addr, _ := runtime.Publish(ctx, PUBLISH_ADDRESS)
	cosign_password := dag.SetSecret("cosign_password", cosignPassword)
	regpassword := dag.SetSecret("reg_password", regPassword)
	_, err := dag.Cosign().Sign(ctx, cosignKey, cosign_password, []string{addr}, dagger.CosignSignOpts{RegistryUsername: regUsername, RegistryPassword: regpassword})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Published to %s 🎉\n", addr)
	return addr
}

func fetchMainGoPath(ctx context.Context, directoryArg *dagger.Directory) (*dagger.Container, string) {

	container := dag.Container().
		From("golang:1.22-alpine").
		WithMountedDirectory("/src", directoryArg).
		WithWorkdir("/src").
		WithExec([]string{"sh", "-c", "export MAIN_GO_PATH=$(find ./cmd -type f -name 'main.go' -print -quit) && echo $MAIN_GO_PATH > main_go_path.txt"})

	// Reading the content of main_go_path.txt file and fetching the actual path of main.go
	main_go_txt_file, _ := container.File("main_go_path.txt").Contents(ctx)
	trimmedPath := strings.TrimPrefix(main_go_txt_file, "./")
	result := "/src/" + trimmedPath
	main_go_path := strings.TrimRight(result, "\n")

	return container, main_go_path
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
