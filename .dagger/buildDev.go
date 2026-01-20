package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"dagger/harbor-cli/internal/dagger"
)

// Create build of Harbor CLI for local testing and development
func (m *HarborCli) BuildDev(ctx context.Context, platform string, source *dagger.Directory) *dagger.File {
	err := m.init(ctx, source)
	if err != nil {
		return nil
	}

	fmt.Println("🛠️  Building Harbor-Cli with Dagger...")
	// Define the path for the binary output
	os, arch, err := parsePlatform(platform)
	if err != nil {
		log.Fatalf("Error parsing platform: %v", err)
	}

	builder := dag.Container().
		From("golang:"+m.GoVersion).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+m.GoVersion)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+m.GoVersion)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithMountedDirectory("/src", m.Source). // Ensure the source directory with go.mod is mounted
		WithWorkdir("/src").
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch)

	gitCommit, _ := builder.WithExec([]string{"git", "rev-parse", "--short", "HEAD", "--always"}).Stdout(ctx)
	buildTime := time.Now().UTC().Format(time.RFC3339)

	ldflagsArgs := LDFlags(ctx, m.AppVersion, m.GoVersion, buildTime, gitCommit)

	builder = builder.WithExec([]string{
		"go", "build", "-ldflags", ldflagsArgs, "-o", "/bin/harbor-cli", "/src/cmd/harbor/main.go",
	})
	return builder.File("/bin/harbor-cli")
}

// Parse the platform string into os and arch
func parsePlatform(platform string) (string, string, error) {
	parts := strings.Split(platform, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid platform format: %s. Should be os/arch. E.g. darwin/amd64", platform)
	}

	return parts[0], parts[1], nil
}
