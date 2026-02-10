package main

import (
	"context"
	"fmt"

	"dagger/harbor-cli/internal/dagger"
)

// Generates Software Bill of Materials for the archive files
func (m *HarborCli) SBOM(ctx context.Context,
	buildDir *dagger.Directory,
	// +ignore=[".gitignore"]
	// +defaultPath="."
	source *dagger.Directory,
) (*dagger.Directory, error) {
	if !m.IsInitialized {
		err := m.init(ctx, source)
		if err != nil {
			return nil, err
		}
	}

	entries, err := buildDir.Entries(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not read dist directory: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("dist directory is empty — run build first")
	}

	goos := []string{"linux", "darwin", "windows"}
	goarch := []string{"amd64", "arm64"}

	sbomFiles := dag.Directory()
	sbomCtr := dag.Container().
		From("alpine:3.19").
		WithExec([]string{
			"sh", "-c",
			`
    set -e
    apk add --no-cache curl ca-certificates tar gzip unzip
    curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh
    syft --version
    `,
		}).
		WithMountedDirectory("/input", buildDir.Directory("archive")).
		WithMountedDirectory("/out", sbomFiles)

	for _, os := range goos {
		for _, arch := range goarch {
			archiveName := fmt.Sprintf("harbor-cli_%s_%s_%s", m.AppVersion, os, arch)
			if os == "windows" {
				archiveName += ".zip"
			} else {
				archiveName += ".tar.gz"
			}

			cmd := []string{
				"sh", "-c",
				fmt.Sprintf(
					"syft /input/%s -o cyclonedx-json > /out/%s.sbom.json",
					archiveName,
					archiveName,
				),
			}

			sbomCtr = sbomCtr.WithExec(cmd)

			sbomFiles = sbomFiles.WithFile(fmt.Sprintf("%s.sbom.json", archiveName), sbomCtr.File(fmt.Sprintf("/out/%s.sbom.json", archiveName)))
		}
	}

	buildDir = buildDir.WithDirectory("sbom", sbomFiles)

	return buildDir, nil
}
