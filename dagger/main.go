package main

import (
	"context"
	"fmt"
	"strings"

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

func (m *HarborCli) Build(ctx context.Context, directoryArg *Directory) {
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
				WithExec([]string{"go", "build", "-o", path + "harbor", main_go_path})
			// Get reference to build output directory in container
			outputs = outputs.WithDirectory(path, build.Directory(path))

		}
	}
	
	_, err := outputs.Export(ctx, "./test")
	if err != nil {
		panic(err)
	}
	fmt.Println("BUILD COMPLETED!✅")
}
// func (m *HarborCli) BuildEnv(directoryArg *Directory) *Container {
// 	return dag.Container().
// 		From("golang:latest").
// 		WithMountedDirectory("/src", directoryArg).
// 		WithWorkdir("/src").
// 		WithExec([]string{"sh", "-c", "export MAIN_GO_PATH=$(find . -type f -name 'main.go' -print -quit) && echo $MAIN_GO_PATH > main_go_path.txt"})
// }

// func getMainGoPath(ctx context.Context, golangcont *Container) string {
// 	// Reading the content of main_go_path.txt file and fetching the actual path of main.go
// 	main_go_txt_file, _ := golangcont.File("main_go_path.txt").Contents(ctx)
// 	trimmedPath := strings.TrimPrefix(main_go_txt_file, "./")
// 	result := "/src/" + trimmedPath
// 	return strings.TrimRight(result, "\n")
// }

// func main() {
// 	ctx := context.Background()
// 	harborCli := &HarborCli{}
// 	directoryArg := dag.Directory()

// 	// Lint code
// 	lintContainer := harborCli.LintCode(ctx, directoryArg)
// 	fmt.Println("Linting completed")

// 	// Build code
// 	buildContainer := harborCli.Build(ctx, directoryArg)
// 	_, err := buildContainer.Export(ctx, ".")
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println("BUILD COMPLETED!✅")
// }
