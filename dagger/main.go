package main

import (
	"context"
	"fmt"
	"log"
	"strings"
)

const (
	GO_VERSION         = "1.22.5"
	SYFT_VERSION       = "v1.9.0"
	GORELEASER_VERSION = "v2.1.0"
	APP_NAME           = "dagger-harbor-cli"
)

type HarborCli struct{}

// example usage: "dagger call container-echo --string-arg yo stdout"
func (m *HarborCli) ContainerEcho(stringArg string) *Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})
}

// example usage: "dagger call grep-dir --directory-arg . --pattern GrepDir"
func (m *HarborCli) GrepDir(ctx context.Context, directoryArg *Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}

func (m *HarborCli) LintCode(ctx context.Context, directoryArg *Directory) *Container {
	fmt.Println("👀 Running linter with Dagger...")
	return dag.Container().
		From("golangci/golangci-lint:latest").
		WithMountedDirectory("/src", directoryArg).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "--timeout", "5m"})
}

func (m *HarborCli) BuildHarbor(ctx context.Context, directoryArg *Directory) *Directory{
	fmt.Println("🛠️  Building with Dagger...")
	oses := []string{"linux", "darwin", "windows"}
	arches := []string{"amd64", "arm64"}

	outputs := dag.Directory()
	golangcont := dag.Container().
		From("golang:latest").
		WithMountedDirectory("/src", directoryArg).
		WithWorkdir("/src").
		WithExec([]string{"sh", "-c", "export MAIN_GO_PATH=$(find . -type f -name 'main.go' -print -quit) && echo $MAIN_GO_PATH > main_go_path.txt"})

	// Reading the content of main_go_path.txt file and fetching the actual path of main.go
	main_go_txt_file, _ := golangcont.File("main_go_path.txt").Contents(ctx)
	trimmedPath := strings.TrimPrefix(main_go_txt_file, "./")
	result := "/src/" + trimmedPath
	main_go_path := strings.TrimRight(result, "\n")

	for _, goos := range oses {
		for _, goarch := range arches {
			path := fmt.Sprintf("build/%s/%s/", goos, goarch)
			build := golangcont.WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch).
				WithExec([]string{"go", "build", "-o", path+"harbor", main_go_path})
			// Get reference to build output directory in container
			outputs = outputs.WithDirectory(path, build.Directory(path))

		}
	}
	return outputs
}

func (m *HarborCli) PullRequest(ctx context.Context, directoryArg *Directory, githubToken string) {

	goreleaser := goreleaserContainer(directoryArg, githubToken).WithExec([]string{"release", "--snapshot", "--clean"})
	_, err := goreleaser.Stderr(ctx)

	if err != nil {
		log.Printf("❌ Error occured during snapshot release for the recently merged pull-request: %s", err)
		return
	}
	log.Println("Pull-Request tasks completed successfully 🎉")
}

// `example: go run ci/dagger.go release`
func (m *HarborCli) Release(ctx context.Context, directoryArg *Directory, githubToken string) {
	goreleaser := goreleaserContainer(directoryArg, githubToken).WithExec([]string{"--clean"})

	_, err := goreleaser.Stderr(ctx)
	if err != nil {
		log.Printf("Error occured during release: %s", err)
		return
	}

	log.Println("Release tasks completed successfully 🎉")
}

func goreleaserContainer(directoryArg *Directory, githubToken string) *Container {
	token := dag.SetSecret("github_token", githubToken)

	// Export the syft binary from the syft container as a file to generate SBOM
	syft := dag.Container().From(fmt.Sprintf("anchore/syft:%s", SYFT_VERSION)).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("gomod")).
		File("/syft")

	// Run go build to check if the binary compiles
	return dag.Container().From(fmt.Sprintf("goreleaser/goreleaser:%s", GORELEASER_VERSION)).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("gomod")).
		WithFile("/bin/syft", syft).
		WithMountedDirectory("/src", directoryArg).WithWorkdir("/src").
		WithEnvVariable("TINI_SUBREAPER", "true").
		WithSecretVariable("GITHUB_TOKEN", token)

}