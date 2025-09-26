package main

import (
	"context"
	"fmt"

	"dagger/harbor-cli/internal/dagger"
)

// +dagger.function
func (m *HarborCli) LintReport(ctx context.Context, source *dagger.Directory) (*dagger.File, error) {
	err := m.Init(ctx, source)
	if err != nil {
		return nil, err
	}

	report := "golangci-lint.report"
	return m.lint(ctx).WithExec([]string{
		"golangci-lint", "run", "-v",
		"--output.tab.path=" + report,
		"--issues-exit-code", "0",
	}).File(report), nil
}

// +dagger.function
func (m *HarborCli) Lint(ctx context.Context, source *dagger.Directory) (string, error) {
	err := m.Init(ctx, source)
	if err != nil {
		return "", err
	}

	return m.lint(ctx).WithExec([]string{"golangci-lint", "run"}).Stderr(ctx)
}

func (m *HarborCli) lint(_ context.Context) *dagger.Container {
	fmt.Println("👀 Running linter and printing results to file golangci-lint.txt.")
	linter := dag.Container().
		From("golangci/golangci-lint:latest-alpine").
		WithMountedCache("/lint-cache", dag.CacheVolume("/lint-cache")).
		WithEnvVariable("GOLANGCI_LINT_CACHE", "/lint-cache").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src")
	return linter
}
