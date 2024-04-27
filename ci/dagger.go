package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
)

const (
	GO_VERSION         = "1.22"
	SYFT_VERSION       = "v0.105.0"
	GORELEASER_VERSION = "v1.24.0"
	APP_NAME           = "dagger-harbor-cli"
	BUILD_PATH         = "dist"
)

var (
	err      error
	res      string
	is_local bool
	GithubToken string
)

func main() {
	// Set a global flag when running locally
	flag.BoolVar(&is_local, "local", false, "whether to run locally [global]")
	flag.Parse()
	task := flag.Arg(0)
	GithubToken = flag.Arg(1)

	if len(task) == 0 {
		log.Fatalln("Missing argument. Expected either 'pull-request' or 'release'.")
	}
	if task != "pull-request" && task != "release" {
		log.Fatalln("Invalid argument. Expected either 'pull-request' or 'release'.")
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer func() {
		log.Printf("Closing Dagger client...")
		client.Close()
	}()

	log.Println("Connected to Dagger")
	switch task {
	case "pull-request":
		res, err = pullrequest(ctx, client)
	case "release":
		res, err = release(ctx, client)
	}

	if err != nil {
		panic(fmt.Sprintf("Error %s: %+v\n", task, err))

	}
	log.Println(res)
}

// `example: go run ci/dagger.go [-local,-help] pull-request `
func pullrequest(ctx context.Context, client *dagger.Client) (string, error) {
	directory := client.Host().Directory(".")

	// Create a go container with the source code mounted
	golang := client.Container().
		From(fmt.Sprintf("golang:%s-alpine", GO_VERSION)).
		WithMountedDirectory("/src", directory).WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", client.CacheVolume("gomod")).
		WithEnvVariable("CGO_ENABLED", "0")

	_, err := golang.WithExec([]string{"go", "test", "./..."}).
		Stderr(ctx)

	if err != nil {
		return "", err
	}

	log.Println("Tests passed successfully!")
	goreleaser := goreleaserContainer(ctx, client, directory).WithExec([]string{"release", "--snapshot", "--clean"})
	_, err = goreleaser.Stderr(ctx)

	if err != nil {
		return "", err
	}
	
	if is_local {
		// Retrieve the dist directory from the container
		dist := goreleaser.Directory(BUILD_PATH)

		// Export the dist directory when running locally
		_, err = dist.Export(ctx, BUILD_PATH)
		if err != nil {
			return "", err
		}
		log.Printf("Exported %v to local successfully!", BUILD_PATH)

	}
	return "Pull-Request tasks completed successfully!", nil
}

// `example: go run ci/dagger.go release`
func release(ctx context.Context, client *dagger.Client) (string, error) {
	directory := client.Host().Directory(".")
	goreleaser := goreleaserContainer(ctx, client, directory).WithExec([]string{"--clean"})

	_, err = goreleaser.Stderr(ctx)
	if err != nil {
		return "", err
	}

	return "Release tasks completed successfully!", nil
}

// goreleaserContainer returns a goreleaser container with the syft binary mounted and GITHUB_TOKEN secret set
//
// `example: goreleaserContainer(ctx, client, directory).WithExec([]string{"build"})`
func goreleaserContainer(ctx context.Context, client *dagger.Client, directory *dagger.Directory) *dagger.Container {
	token := client.SetSecret("github_token", GithubToken)

	// Export the syft binary from the syft container as a file
	syft := client.Container().From(fmt.Sprintf("anchore/syft:%s", SYFT_VERSION)).
		WithMountedCache("/go/pkg/mod", client.CacheVolume("gomod")).
		File("/syft")

	// Run go build to check if the binary compiles
	return client.Container().From(fmt.Sprintf("goreleaser/goreleaser:%s", GORELEASER_VERSION)).
		WithMountedCache("/go/pkg/mod", client.CacheVolume("gomod")).
		WithFile("/bin/syft", syft).
		WithMountedDirectory("/src", directory).WithWorkdir("/src").
		WithEnvVariable("TINI_SUBREAPER", "true").
		WithSecretVariable("GITHUB_TOKEN", token)

}