package main

import (
	"context"
	"fmt"
	"strings"
	"os"

	"dagger.io/dagger"
)

func main() {
	if err := build(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func build(ctx context.Context) error {
	fmt.Println("Building with Dagger")
	oses := []string{"linux", "darwin", "windows"}
	arches := []string{"amd64", "arm64"}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()
	src := client.Host().Directory(".")
	outputs := client.Directory()
	golangcont := client.Container().From("golang:latest")
	golangcont = golangcont.WithDirectory("/src", src).WithWorkdir("/src")

	golangcont = golangcont.WithExec([]string{"sh", "-c", "export MAIN_GO_PATH=$(find . -type f -name 'main.go' -print -quit) && echo $MAIN_GO_PATH > main_go_path.txt" })

	// reading the content of main_go_path.txt file and fetching the actual path of main.go
	main_go_txt_file, _ := golangcont.File("main_go_path.txt").Contents(ctx)
	trimmedPath := strings.TrimPrefix(main_go_txt_file, "./")
	result := "/src/" + trimmedPath
	main_go_path := strings.TrimRight(result, "\n")
	
	for _, goos := range oses {
		for _, goarch := range arches {
			path := fmt.Sprintf("build/%s/%s/", goos, goarch)
			build := golangcont.WithEnvVariable("GOOS", goos)
			build = build.WithEnvVariable("GOARCH", goarch)
			build = build.WithExec([]string{"go", "build", "-o", path, main_go_path})
			// get reference to build output directory in container
			outputs = outputs.WithDirectory(path, build.Directory(path))
		}
	}
	// write build artifacts to host
	_, err = outputs.Export(ctx, ".")
	if err != nil {
		return err
	}
	fmt.Println("BUILD COMPLETED!")
	return nil
}
