package pipeline

import (
	"context"
	"fmt"

	"dagger/harbor-cli/internal/dagger"
)

func (s *Pipeline) AptRepoBuild(ctx context.Context, dist *dagger.Directory, token *dagger.Secret) error {
	archs := []string{"amd64", "arm64"}
	root := s.dag.Directory()
	root = root.WithDirectory("pool/main/m", dist.Directory("deb"))
	githubToken, err := token.Plaintext(ctx)
	if err != nil {
		return err
	}

	// Base container
	container := s.dag.Container().
		From("debian:bookworm-slim").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "dpkg-dev", "gzip", "git"}).
		WithEnvVariable("GH_TOKEN", githubToken).
		WithMountedDirectory("/repo", root)

	// Building `Package` file for each arch
	for _, arch := range archs {
		pkgDir := fmt.Sprintf("/repo/dists/stable/main/binary-%s", arch)
		poolDir := "/repo/pool/main/m"

		container = container.WithExec([]string{
			"bash", "-c",
			fmt.Sprintf("mkdir -p %s && dpkg-scanpackages -a %s %s /dev/null > %s/Packages && gzip -9c %s/Packages > %s/Packages.gz",
				pkgDir, arch, poolDir, pkgDir, pkgDir, pkgDir),
		})
	}

	container = container.WithExec([]string{
		"bash", "-c",
		`cat <<EOF > /repo/dists/stable/Release
Origin: https://github.com/nucleofusion/harbor-cli  
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
        git remote add origin https://x-access-token:$GH_TOKEN@github.com/nucleofusion/harbor-cli.git
        git checkout -B gh-pages || git checkout --orphan gh-pages

        git config user.name "github-actions[bot]"
        git config user.email "github-actions[bot]@users.noreply.github.com"

        git add dists pool 

        git commit -m "Update APT repo for %s" || echo "No changes to commit"
        git push origin gh-pages -f
        `, s.appVersion),
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
