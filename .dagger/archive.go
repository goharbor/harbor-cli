package main

import (
	"context"
	"fmt"

	"dagger/harbor-cli/internal/dagger"
)

func (m *HarborCli) Archive(ctx context.Context,
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

	archives := dag.Directory()

	for _, os := range goos {
		for _, arch := range goarch {
			binName := fmt.Sprintf("harbor-cli_%s_%s_%s", m.AppVersion, os, arch)
			if os == "windows" {
				binName += ".exe"
			}

			binPath := fmt.Sprintf("bin/%s", binName)

			archiveName := fmt.Sprintf("harbor-cli_%s_%s_%s", m.AppVersion, os, arch)
			archiveDir := getArchiveDirectory(buildDir.File(binPath), source)

			var (
				archiveFile string
				container   *dagger.Container
			)

			if os == "windows" {
				// Handle Windows .zip
				archiveFile = archiveName + ".zip"
				container = dag.Container().
					From("alpine:latest").
					WithExec([]string{"apk", "add", "--no-cache", "zip"}).
					WithMountedDirectory("/input", archiveDir).
					WithMountedDirectory("/out", archives).
					WithWorkdir("/input").
					WithExec([]string{"zip", "-j", "/out/" + archiveFile, "/input/harbor-cli", "/input/LICENSE", "/input/README.md"})
			} else {
				archiveFile = archiveName + ".tar.gz"
				container = dag.Container().
					From("alpine:latest").
					WithMountedDirectory("/input", archiveDir).
					WithMountedDirectory("/out", archives).
					WithWorkdir("/input").
					WithExec([]string{
						"tar", "-czf", "/out/" + archiveFile, "-C", "/input", ".",
					})
			}

			archives = archives.WithFile(archiveFile, container.File("/out/"+archiveFile))
		}
	}

	buildDir = buildDir.WithDirectory("archive", archives)

	return buildDir, nil
}

func getArchiveDirectory(harborCli *dagger.File, source *dagger.Directory) *dagger.Directory {
	archiveDir := dag.Directory()

	archiveDir = archiveDir.
		WithFile("LICENSE", source.File("LICENSE")).
		WithFile("README.md", source.File("README.md")).
		WithFile("harbor-cli", harborCli)

	return archiveDir
}
