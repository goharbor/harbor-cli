package main

import (
	"context"
	"fmt"
	"regexp"
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
	GoVersion  string
}

// +dagger.function
func (m *HarborCli) Pipeline(ctx context.Context, source *dagger.Directory, githubToken *dagger.Secret) (*dagger.Directory, error) {
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

	// Building Brew Formula
	dist, err = pipe.BrewFormula(ctx, dist)
	if err != nil {
		return nil, err
	}

	// Publishing Release
	out, err := pipe.PublishRelease(ctx, dist, githubToken)
	if err != nil {
		return nil, err
	}
	fmt.Println(out)

	// Publishing repo
	err = pipe.AptRepoBuild(ctx, dist, githubToken)
	if err != nil {
		return nil, err
	}

	return dist, err
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
