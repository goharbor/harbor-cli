package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/harbor-cli/internal/dagger"
	"dagger/harbor-cli/pipeline"
)

const (
	GO_VERSION           = "1.24.2"
	GOLANGCILINT_VERSION = "1.24.2"
)

type HarborCli struct {
	Source     *dagger.Directory
	AppVersion string
}

// +dagger.function
func (m *HarborCli) Pipeline(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	err := m.Init(ctx, source)
	if err != nil {
		return nil, err
	}

	dist := dag.Directory()
	pipe := pipeline.InitPipeline(source, dag, m.AppVersion)

	// Building Binaries
	dist, err = pipe.Build(ctx, dist, GO_VERSION)
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

	return dist, nil
}

func (m *HarborCli) Init(ctx context.Context, source *dagger.Directory) error {
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

	m.Source = source
	m.AppVersion = strings.TrimSpace(out)

	return nil
}
