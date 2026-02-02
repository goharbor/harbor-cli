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
		WithExec([]string{"apt-get", "install", "-y", "dpkg-dev", "gzip", "git"}).
		WithEnvVariable("GH_TOKEN", githubToken).
		WithMountedDirectory("/repo", root).
		WithWorkdir("/repo")

	// Building `Package` file for each arch
	for _, arch := range archs {
		pkgDir := fmt.Sprintf("buildDirs/stable/main/binary-%s", arch)
		poolDir := "pool/main/m"

		container = container.WithExec([]string{
			"bash", "-c",
			fmt.Sprintf("mkdir -p %s && dpkg-scanpackages -a %s %s /dev/null > %s/Packages && gzip -9c %s/Packages > %s/Packages.gz && rm -rf %s/Packages",
				pkgDir, arch, poolDir, pkgDir, pkgDir, pkgDir, pkgDir),
		})
	}

	// Release File
	container = container.WithExec([]string{
		"bash", "-c",
		`cat <<EOF > /repo/buildDirs/stable/Release
Origin: https://github.com/goharbor/harbor-cli  
Label: HarborCLI 
Suite: stable
Codename: stable
Architectures: amd64 arm64
Components: main
Description: Harbor CLI — a command-line interface for interacting with your Harbor container registry.
EOF`,
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

        git add buildDirs pool 

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
// ├── dist/
// │   └── stable/
// │       ├── Release
// │       └── main/
// │           ├── binary-amd64/
// │           │   └── Packages.gz
// │           └── binary-arm64/
// │               └── Packages.gz
// └── pool/
//     └── main/
//         └── m/
//             ├── myapp_1.0.0_amd64.deb
//             └── myapp_1.0.0_arm64.deb
//
