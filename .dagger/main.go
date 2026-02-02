package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"dagger/harbor-cli/internal/dagger"
)

type HarborCli struct {
	Source        *dagger.Directory // Source Directory where code resides
	AppVersion    string            // Current Version of the app, acquired from git tags
	GoVersion     string            // Go Version used in the current release, acquired from the go.mod file
	IsInitialized bool
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
	m.AppVersion = strings.TrimSpace(strings.TrimLeft(out, "v"))
	m.IsInitialized = true

	return nil
}
