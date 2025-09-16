package pipeline

import (
	"context"
	"fmt"

	"dagger/harbor-cli/internal/dagger"
)

func (s *Pipeline) Archive(ctx context.Context, dist *dagger.Directory) (*dagger.Directory, error) {
	entries, err := dist.Entries(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not read dist directory: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("dist directory is empty — run build first")
	}

	goos := []string{"linux", "darwin", "windows"}
	goarch := []string{"amd64", "arm64"}

	archives := s.dag.Directory()

	for _, os := range goos {
		for _, arch := range goarch {
			binName := fmt.Sprintf("harbor-cli_%s_%s_%s", s.appVersion, os, arch)
			if os == "windows" {
				binName += ".exe"
			}

			binPath := fmt.Sprintf("%s/%s", os, binName)

			archiveName := fmt.Sprintf("harbor-cli_%s_%s_%s", s.appVersion, os, arch)

			var (
				archiveFile string
				container   *dagger.Container
			)

			if os == "windows" {
				// Handle Windows .zip
				archiveFile = archiveName + ".zip"
				container = s.dag.Container().
					From("alpine:latest").
					WithExec([]string{"apk", "add", "--no-cache", "zip"}).
					WithMountedDirectory("/input", dist).
					WithMountedDirectory("/out", archives).
					WithWorkdir("/input").
					WithExec([]string{"zip", "-j", "/out/" + archiveFile, binPath})
			} else {
				archiveFile = archiveName + ".tar.gz"
				container = s.dag.Container().
					From("alpine:latest").
					WithMountedDirectory("/input", dist).
					WithMountedDirectory("/out", archives).
					WithWorkdir("/input").
					WithExec([]string{
						"tar", "-czf", "/out/" + archiveFile, "-C", os + "/", binName,
					})
			}

			archives = archives.WithFile(archiveFile, container.File("/out/"+archiveFile))
		}
	}

	dist = dist.WithDirectory("archive", archives)

	return dist, nil
}
