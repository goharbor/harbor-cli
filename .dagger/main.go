package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"dagger/harbor-cli/internal/dagger"
	"dagger/harbor-cli/pipeline"
)

type HarborCli struct {
	Source     *dagger.Directory // Source Directory where code resides
	AppVersion string            // Current Version of the app, acquired from git tags
	GoVersion  string            // Go Version used in the current release, acquired from the go.mod file
}

// The _full_ pipeline for CI/CD
// Build Binaries -> Generate zip/tar.gz -> Building .deb & .rpm -> Building Brew Formula
// -> Publishing to release page -> Publishing to apt
func (m *HarborCli) Pipeline(ctx context.Context, source *dagger.Directory,
	// Secrets
	githubToken *dagger.Secret, registryPassword *dagger.Secret,
	registryAddr string, registryUsername string,
) (*dagger.Directory, error) {
	err := m.init(ctx, source)
	if err != nil {
		return nil, err
	}

	dist := dag.Directory()
	pipe := pipeline.InitPipeline(source, dag, m.AppVersion, m.GoVersion)

	// Building Binaries
	dist, err = pipe.Build(ctx, dist)
	if err != nil {
		return nil, err
	}

	// Archiving Binaries
	dist, err = pipe.Archive(ctx, dist)
	if err != nil {
		return nil, err
	}

	// Building deb/rpm Binaries
	dist, err = pipe.NFPMBuild(ctx, dist)
	if err != nil {
		return nil, err
	}

	// Building Checksum file
	dist, err = pipe.Checksum(ctx, dist)
	if err != nil {
		return nil, err
	}

	// Publishing Release
	out, err := pipe.PublishRelease(ctx, dist, githubToken)
	if err != nil {
		return nil, err
	}
	fmt.Println(out)

	// Publishing Image
	res := pipe.PublishImage(ctx, dist, registryAddr, registryUsername, []string{"latest", m.AppVersion}, registryPassword)
	if err != nil {
		return nil, err
	}
	fmt.Println(strings.Join(res, "\n"))

	// Publishing repo
	err = pipe.AptRepoBuild(ctx, dist, githubToken)
	if err != nil {
		return nil, err
	}

	return dist, err
}

func (m *HarborCli) init(ctx context.Context, source *dagger.Directory) error {
	out, err := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"git", "describe", "--tags", "--abbrev=0"}).
		Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	goVersion, err := source.File("go.mod").Contents(ctx)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`(?m)^go (\d+\.\d+(\.\d+)?)`)
	match := re.FindStringSubmatch(goVersion)
	if len(match) > 1 {
		m.GoVersion = match[1]
	}

	m.Source = source
	m.AppVersion = strings.TrimSpace(out)

	return nil
}
