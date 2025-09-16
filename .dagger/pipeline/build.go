package pipeline

import (
	"context"
	"fmt"

	"dagger/harbor-cli/internal/dagger"
)

func (s *Pipeline) Build(ctx context.Context, dist *dagger.Directory, GO_VERSION string) (*dagger.Directory, error) {
	goos := []string{"linux", "darwin", "windows"}
	goarch := []string{"amd64", "arm64"}

	for _, os := range goos {
		for _, arch := range goarch {
			// Defining binary file name
			binName := fmt.Sprintf("harbor-cli_%s_%s_%s", s.appVersion, os, arch)
			if os == "windows" {
				binName += ".exe"
			}

			builder := s.dag.Container().
				From("golang:"+GO_VERSION).
				WithMountedCache("/go/pkg/mod", s.dag.CacheVolume("go-mod-"+GO_VERSION)).
				WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
				WithMountedCache("/go/build-cache", s.dag.CacheVolume("go-build-"+GO_VERSION)).
				WithEnvVariable("GOCACHE", "/go/build-cache").
				WithMountedDirectory("/src", s.source).
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
