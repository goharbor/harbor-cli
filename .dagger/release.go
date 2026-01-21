package main

import (
	"context"

	"dagger/harbor-cli/internal/dagger"
)

func (m *HarborCli) PublishRelease(ctx context.Context,
	buildDir *dagger.Directory,
	// +ignore=[".gitignore"]
	// +defaultPath="."
	source *dagger.Directory,
	token *dagger.Secret,
) (string, error) {
	if !m.IsInitialized {
		err := m.init(ctx, source)
		if err != nil {
			return "", err
		}
	}

	bins, err := DistBinaries(ctx, dag, buildDir)
	if err != nil {
		return "", err
	}

	cmd := []string{"gh", "release", "upload", "v" + m.AppVersion}
	cmd = append(cmd, bins...)
	cmd = append(cmd, "/dist/checksums.txt")
	cmd = append(cmd, "--clobber")

	ctr := dag.Container().
		From("debian:bookworm-slim").
		WithMountedDirectory("/src", source).
		WithMountedDirectory("/dist", buildDir).
		WithSecretVariable("GH_TOKEN", token).
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "curl", "git"}).
		WithExec([]string{"curl", "-fsSL", "https://cli.github.com/packages/githubcli-archive-keyring.gpg", "-o", "/usr/share/keyrings/githubcli-archive-keyring.gpg"}).
		WithExec([]string{"sh", "-c", `echo "deb [arch=amd64 signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" > /etc/apt/sources.list.d/github-cli.list`}).
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "gh"})

	return ctr.
		WithWorkdir("/src").
		// Create release if it doesn't exist, otherwise continue
		WithExec([]string{"sh", "-c", "gh release view v" + m.AppVersion + " || gh release create v" + m.AppVersion + " --generate-notes"}).
		WithExec(cmd).
		Stdout(ctx)
}
