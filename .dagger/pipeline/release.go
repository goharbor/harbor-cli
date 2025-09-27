package pipeline

import (
	"context"

	"dagger/harbor-cli/internal/dagger"
)

func (s *Pipeline) PublishRelease(ctx context.Context, dist *dagger.Directory) (string, error) {
	bins, err := DistBinaries(ctx, dist)
	if err != nil {
		return "", err
	}

	cmd := []string{
		"gh", "release", "upload", s.appVersion,
	}
	cmd = append(cmd, bins...)
	cmd = append(cmd, "--clobber")

	return s.dag.Container().
		From("ghcr.io/cli/cli:latest").
		WithMountedDirectory("/dist", dist).
		WithSecretVariable("GH_TOKEN", s.GithubToken).
		WithExec(cmd).Stderr(ctx)
}

func DistBinaries(ctx context.Context, dist *dagger.Directory) ([]string, error) {
	dirs := []string{"archive", "linux", "windows", "darwin", "deb", "rpm"}
	var files []string

	for _, d := range dirs {
		subdir := dist.Directory(d)
		entries, err := subdir.Entries(ctx)
		if err != nil {
			// skip missing directories or return error
			continue
		}

		files = append(files, entries...)
	}

	return files, nil
}
