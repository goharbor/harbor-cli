package main

import (
	"context"
	"fmt"
	"time"

	"dagger/harbor-cli/internal/dagger"
)

func (m *HarborCli) Build(ctx context.Context,
	buildDir *dagger.Directory,
	// +ignore=[".gitignore"]
	// +defaultPath="."
	source *dagger.Directory,
) (*dagger.Directory, error) {
	if !m.IsInitialized {
		err := m.init(ctx, source)
		if err != nil {
			return nil, err
		}
	}

	goos := []string{"linux", "darwin", "windows"}
	goarch := []string{"amd64", "arm64"}

	for _, os := range goos {
		for _, arch := range goarch {
			// Defining binary file name
			binName := fmt.Sprintf("harbor-cli_%s_%s_%s", m.AppVersion, os, arch)
			if os == "windows" {
				binName += ".exe"
			}

			builder := dag.Container().
				From("golang:"+m.GoVersion+"-alpine").
				WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+m.GoVersion)).
				WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
				WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+m.GoVersion)).
				WithEnvVariable("GOCACHE", "/go/build-cache").
				WithMountedDirectory("/src", source).
				WithWorkdir("/src").
				WithEnvVariable("GOOS", os).
				WithEnvVariable("GOARCH", arch)

			gitCommit, _ := builder.WithExec([]string{"git", "rev-parse", "--short", "HEAD", "--always"}).Stdout(ctx)
			buildTime := time.Now().UTC().Format(time.RFC3339)

			ldflagsArgs := LDFlags(ctx, m.AppVersion, m.GoVersion, buildTime, gitCommit)

			builder = builder.WithExec([]string{
				"sh", "-c",
				fmt.Sprintf(`go build -v -ldflags "%s" -o /bin/%s /src/cmd/harbor/main.go`, ldflagsArgs, binName),
			})

			file := builder.File("/bin/" + binName)                             // Taking file from container
			buildDir = buildDir.WithFile(fmt.Sprintf("/bin/%s", binName), file) // Adding file(bin) to dist directory
		}
	}

	return buildDir, nil
}
