package main

import (
	"context"
	"fmt"

	"dagger/harbor-cli/internal/dagger"
)

// Checks for vulnerabilities using govulncheck
func (m *HarborCli) vulnerabilityCheck(ctx context.Context) *dagger.Container {
	return dag.Container().
		From("golang:"+m.GoVersion+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+m.GoVersion)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+m.GoVersion)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithExec([]string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@latest"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src")
}

// Runs a vulnerability check using govulncheck
func (m *HarborCli) VulnerabilityCheck(ctx context.Context, source *dagger.Directory) (string, error) {
	err := m.init(ctx, source)
	if err != nil {
		return "", err
	}
	return m.vulnerabilityCheck(ctx).
		WithExec([]string{"govulncheck", "-show", "verbose", "./..."}).
		Stderr(ctx)
}

// Runs a vulnerability check using govulncheck and writes results to vulnerability-check.report
func (m *HarborCli) VulnerabilityCheckReport(ctx context.Context, source *dagger.Directory) (*dagger.File, error) {
	err := m.init(ctx, source)
	if err != nil {
		return nil, err
	}

	report := "vulnerability-check.report"
	cmd := fmt.Sprintf("govulncheck ./... > %s || true", report)

	return m.vulnerabilityCheck(ctx).
		WithExec([]string{
			"sh", "-c", cmd,
		}).File(report), nil
}
