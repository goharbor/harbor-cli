package pipeline

import (
	"context"
	"fmt"

	"dagger/harbor-cli/internal/dagger"
)

func (s *Pipeline) PublishRelease(ctx context.Context, dist *dagger.Directory, token *dagger.Secret) (string, error) {
	bins, err := DistBinaries(ctx, s.dag, dist)
	if err != nil {
		return "", err
	}

	cmd := []string{"gh", "release", "upload", s.appVersion}
	cmd = append(cmd, bins...)
	cmd = append(cmd, "/dist/checksum.txt")
	cmd = append(cmd, "--clobber")

	ctr := s.dag.Container().
		From("debian:bookworm-slim").
		WithMountedDirectory("/src", s.source).
		WithMountedDirectory("/dist", dist).
		WithSecretVariable("GH_TOKEN", token).
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "curl", "git"}).
		WithExec([]string{"curl", "-fsSL", "https://cli.github.com/packages/githubcli-archive-keyring.gpg", "-o", "/usr/share/keyrings/githubcli-archive-keyring.gpg"}).
		WithExec([]string{"sh", "-c", `echo "deb [arch=amd64 signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" > /etc/apt/sources.list.d/github-cli.list`}).
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "gh"})

	return ctr.
		WithWorkdir("/src").
		// Creating Release
		WithExec([]string{"gh", "release", "create", s.appVersion, "--title", fmt.Sprintf("Release %s", s.appVersion)}).
		WithExec(cmd).
		Stdout(ctx)
}
