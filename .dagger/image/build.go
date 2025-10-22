package image

import (
	"context"
	"fmt"
	"time"

	"dagger/harbor-cli/internal/dagger"
	"dagger/harbor-cli/utils"
)

func (s *ImagePipeline) Build(ctx context.Context, dist *dagger.Directory) (*dagger.Directory, error) {
	goarch := []string{"amd64", "arm64"}

	for _, arch := range goarch {
		// Defining binary file name
		binName := fmt.Sprintf("harbor-cli_%s_%s_%s", s.appVersion, "linux", arch)

		builder := s.dag.Container().
			From("golang:"+s.goVersion).
			WithMountedCache("/go/pkg/mod", s.dag.CacheVolume("go-mod-"+s.goVersion)).
			WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
			WithMountedCache("/go/build-cache", s.dag.CacheVolume("go-build-"+s.goVersion)).
			WithEnvVariable("GOCACHE", "/go/build-cache").
			WithMountedDirectory("/src", s.source).
			WithWorkdir("/src").
			WithEnvVariable("GOOS", "linux").
			WithEnvVariable("GOARCH", arch).
			WithEnvVariable("CGO_ENABLED", "0")

		gitCommit, _ := builder.WithExec([]string{"git", "rev-parse", "--short", "HEAD", "--always"}).Stdout(ctx)
		buildTime := time.Now().UTC().Format(time.RFC3339)

		ldflagsArgs := utils.LDFlags(ctx, s.appVersion, s.goVersion, buildTime, gitCommit)

		builder = builder.WithExec([]string{
			"bash", "-c",
			fmt.Sprintf(`set -ex && go env && go build -v -ldflags "%s" -o /bin/%s /src/cmd/harbor/main.go`, ldflagsArgs, binName),
		})

		file := builder.File("/bin/" + binName)                            // Taking file from container
		dist = dist.WithFile(fmt.Sprintf("%s/%s", "linux", binName), file) // Adding file(bin) to dist directory
	}

	return dist, nil
}
