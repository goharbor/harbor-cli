package main

import (
	"context"
	"fmt"
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
	arches := []string{"amd64", "arm64", "arm"}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()
	src := client.Host().Directory(".")
	outputs := client.Directory()
	golang := client.Container().From("golang:latest")
	golang = golang.WithDirectory("/src", src).WithWorkdir("/src")

	for _, goos := range oses {
		for _, goarch := range arches {
			path := fmt.Sprintf("build/%s/%s/", goos, goarch)
			build := golang.WithEnvVariable("GOOS", goos)
			build = build.WithEnvVariable("GOARCH", goarch)
			build = build.WithExec([]string{"go", "build", "-o", path})
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