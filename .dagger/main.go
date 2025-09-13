package main

import (
	"context"

	"dagger/harbor-cli/internal/dagger"
)

const (
	GO_VERSION           = "1.24.2"
	GOLANGCILINT_VERSION = "1.24.2"
)

type HarborCli struct {
	Source *dagger.Directory
}

// +dagger.function
func (m *HarborCli) Init(ctx context.Context, source *dagger.Directory) (*HarborCli, error) {
	return &HarborCli{
		Source: source,
	}, nil
}
