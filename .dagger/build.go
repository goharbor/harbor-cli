package main

import (
	"context"
	"fmt"

	"dagger/harbor-cli/internal/dagger"
)

// +dagger.function
func (s *HarborCli) Build(ctx context.Context) (*dagger.Directory, error) {
	goos := []string{"linux", "darwin", "windows"}
	goarch := []string{"amd64", "arm64"}

	// Where all the binaries etc will reside
	dist := dag.Directory()

	for _, os := range goos {
		for _, arch := range goarch {
			// Defining binary file name
			binName := fmt.Sprintf("harbor-cli_%s_%s", os, arch)
			if os == "windows" {
				binName += ".exe"
			}

			builder := dag.Container().
				From("golang:"+GO_VERSION).
				WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+GO_VERSION)).
				WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
				WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+GO_VERSION)).
				WithEnvVariable("GOCACHE", "/go/build-cache").
				WithMountedDirectory("/src", s.Source).
				WithWorkdir("/src").
				WithEnvVariable("GOOS", os).
				WithEnvVariable("GOARCH", arch).
				WithExec([]string{
					"go", "build", "-o", "/bin/" + binName, "./cmd/harbor",
				})

			file := builder.File("/bin/" + binName)                       // Taking file from container
			dist = dist.WithFile(fmt.Sprintf("%s/%s", os, binName), file) // Adding file(bin) to dist directory
		}
	}

	return dist, nil
}

// +dagger.function
func (s *HarborCli) Archive(ctx context.Context) (*dagger.Directory, error) {
	dist, err := s.Build(ctx)
	if err != nil {
		return nil, err
	}

	entries, err := dist.Entries(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not read dist directory: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("dist directory is empty — run build first")
	}

	goos := []string{"linux", "darwin", "windows"}
	goarch := []string{"amd64", "arm64"}

	archives := dag.Directory()

	for _, os := range goos {
		for _, arch := range goarch {
			binName := fmt.Sprintf("harbor-cli_%s_%s", os, arch)
			if os == "windows" {
				binName += ".exe"
			}

			binPath := fmt.Sprintf("%s/%s", os, binName)

			archiveName := fmt.Sprintf("harbor-cli_%s_%s", os, arch)

			var (
				archiveFile string
				container   *dagger.Container
			)

			if os == "windows" {
				// Handle Windows .zip
				archiveFile = archiveName + ".zip"
				container = dag.Container().
					From("alpine:latest").
					WithExec([]string{"apk", "add", "--no-cache", "zip"}).
					WithMountedDirectory("/input", dist).
					WithMountedDirectory("/out", archives).
					WithWorkdir("/input").
					WithExec([]string{"zip", "-j", "/out/" + archiveFile, binPath})
			} else {
				archiveFile = archiveName + ".tar.gz"
				container = dag.Container().
					From("alpine:latest").
					WithMountedDirectory("/input", dist).
					WithMountedDirectory("/out", archives).
					WithWorkdir("/input").
					WithExec([]string{
						"tar", "-czf", "/out/" + archiveFile, "-C", os + "/", binName,
					})
			}

			archives = archives.WithFile(archiveFile, container.File("/out/"+archiveFile))
		}
	}

	dist = dist.WithDirectory("archive", archives)

	return dist, nil
}
