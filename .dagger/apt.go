package main

import (
	"context"
	"fmt"

	"dagger/harbor-cli/internal/dagger"
)

func (m *HarborCli) AptBuild(ctx context.Context,
	buildDir *dagger.Directory,
	// +ignore=[".gitignore"]
	// +defaultPath="."
	source *dagger.Directory,
	token *dagger.Secret,
) error {
	if !m.IsInitialized {
		err := m.init(ctx, source)
		if err != nil {
			return err
		}
	}

	archs := []string{"amd64", "arm64"}
	root := dag.Directory()
	root = root.WithDirectory("pool/main/m", buildDir.Directory("deb"))
	githubToken, err := token.Plaintext(ctx)
	if err != nil {
		return err
	}

	// Base container
	container := dag.Container().
		From("debian:bookworm-slim").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "dpkg-dev", "apt-utils", "gzip", "git"}).
		WithEnvVariable("GH_TOKEN", githubToken).
		WithMountedDirectory("/repo", root).
		WithWorkdir("/repo")

	// Building `Package` file for each arch
	for _, arch := range archs {
		pkgDir := fmt.Sprintf("dists/stable/main/binary-%s", arch)
		poolDir := "pool/main/m"

		container = container.WithExec([]string{
			"bash", "-c",
			fmt.Sprintf("mkdir -p %s && dpkg-scanpackages -a %s %s /dev/null > %s/Packages && gzip -9kf %s/Packages",
				pkgDir, arch, poolDir, pkgDir, pkgDir),
		})
	}

	// Release file. apt-ftparchive adds the Date and the checksums of every
	// Packages index, so apt stops warning about a Release file without hashes.
	container = container.WithExec([]string{
		"bash", "-c",
		`apt-ftparchive \
  -o APT::FTPArchive::Release::Origin="https://github.com/goharbor/harbor-cli" \
  -o APT::FTPArchive::Release::Label="HarborCLI" \
  -o APT::FTPArchive::Release::Suite="stable" \
  -o APT::FTPArchive::Release::Codename="stable" \
  -o APT::FTPArchive::Release::Architectures="amd64 arm64" \
  -o APT::FTPArchive::Release::Components="main" \
  -o APT::FTPArchive::Release::Description="Harbor CLI, a command-line interface for the Harbor container registry" \
  release dists/stable > /tmp/Release && mv /tmp/Release dists/stable/Release`,
	})

	container = container.
		WithWorkdir("/repo").
		WithExec([]string{
			"bash", "-c",
			fmt.Sprintf(`
        set -e
        cd /repo

        git init
        git remote add origin https://x-access-token:$GH_TOKEN@github.com/goharbor/harbor-cli.git
        git checkout -B gh-pages || git checkout --orphan gh-pages

        git config user.name "github-actions[bot]"
        git config user.email "github-actions[bot]@users.noreply.github.com"

        git add dists pool

        git commit -m "Update APT repo for %s" || echo "No changes to commit"
        git push origin gh-pages -f
        `, m.AppVersion),
		})

	_, err = container.Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to run container: %w", err)
	}

	return nil
}

// GH-PAGES Structure
//
// /
// ├── dists/
// │   └── stable/
// │       ├── Release
// │       └── main/
// │           ├── binary-amd64/
// │           │   ├── Packages
// │           │   └── Packages.gz
// │           └── binary-arm64/
// │               ├── Packages
// │               └── Packages.gz
// └── pool/
//     └── main/
//         └── m/
//             ├── myapp_1.0.0_amd64.deb
//             └── myapp_1.0.0_arm64.deb
//
