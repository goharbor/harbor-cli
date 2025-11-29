package main

import (
	"context"

	"dagger/harbor-cli/internal/dagger"
)

// Generate CLI Documentation and return the directory containing the generated files
func (m *HarborCli) RunDoc(ctx context.Context, source *dagger.Directory) (*dagger.Directory, error) {
	err := m.init(ctx, source)
	if err != nil {
		return nil, err
	}

	return dag.Container().
		From("golang:"+m.GoVersion+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+m.GoVersion)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+m.GoVersion)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src/doc").
		WithExec([]string{"go", "run", "doc.go"}).
		WithExec([]string{"go", "run", "./man-docs/man_doc.go"}).
		WithWorkdir("/src").Directory("/src/doc"), nil
}
