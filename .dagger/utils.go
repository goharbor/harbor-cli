package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"dagger/harbor-cli/internal/dagger"
)

func DistBinaries(ctx context.Context, s *dagger.Client, dist *dagger.Directory) ([]string, error) {
	dirs := []string{"archive", "deb", "rpm", "apk"}
	var files []string

	ctr := s.Container().
		From("alpine:latest").
		WithMountedDirectory("/dist", dist).
		WithWorkdir("/dist")

	for _, v := range dirs {
		out, err := ctr.WithExec([]string{"ls", v}).Stdout(ctx)
		if err != nil {
			return nil, err
		}

		bins := strings.Split(out, "\n")
		for _, bin := range bins {
			if bin != "" && bin != "nfpm.yml" {
				files = append(files, filepath.Join("/", "dist", v, bin))
			}
		}
	}

	return files, nil
}

func LDFlags(ctx context.Context, version, goVersion, buildTime, commit string) string {
	return fmt.Sprintf("-X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.Version=%s "+
		"-X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.GoVersion=%s "+
		"-X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.BuildTime=%s "+
		"-X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.GitCommit=%s",
		version, goVersion, buildTime, commit,
	)
}
