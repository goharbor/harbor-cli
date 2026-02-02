package main

import (
	"context"
	"fmt"

	"dagger/harbor-cli/internal/dagger"
)

func (m *HarborCli) Apk(ctx context.Context,
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

	buildfile := dag.File("APKBUILD", apkbuild(m.AppVersion))

	archs := []struct {
		Arch    string
		ApkArch string
	}{
		{"arm64", "aarch64"},
		{"amd64", "x86_64"},
	}

	for _, arch := range archs {
		filename := fmt.Sprintf("bin/harbor-cli_%s_linux_%s", m.AppVersion, arch.Arch)
		binary := buildDir.File(filename)

		apk := dag.Container(dagger.ContainerOpts{
			Platform: dagger.Platform(fmt.Sprintf("linux/%s", arch.Arch)),
		}).
			From("alpine:3.19").
			WithExec([]string{"apk", "add", "--no-cache", "alpine-sdk", "abuild"}).
			WithWorkdir("/build").
			WithFile("/build/harbor-cli", binary).
			WithFile("/build/APKBUILD", buildfile).

			// create builder user + abuild group
			WithExec([]string{"adduser", "-D", "builder"}).
			WithExec([]string{"addgroup", "builder", "abuild"}).
			WithExec([]string{"chown", "-R", "builder:builder", "/build"}).

			// switch to builder FIRST
			WithUser("builder").
			WithEnvVariable("HOME", "/home/builder").

			// generate signing key AS BUILDER (this is critical)
			WithExec([]string{"abuild-keygen", "-a", "-n"}).

			// switch back to root ONLY to trust the public key
			WithUser("root").
			WithExec([]string{
				"sh", "-c",
				"mkdir -p /etc/apk/keys && cp /home/builder/.abuild/*.rsa.pub /etc/apk/keys/",
			}).

			// back to builder for the build
			WithUser("builder").
			WithEnvVariable("HOME", "/home/builder").

			// sanity check (keep this until stable)
			WithExec([]string{"sh", "-c", "ls -l ~/.abuild && echo HOME=$HOME"}).

			// run abuild
			WithExec([]string{"abuild", "-rd"})

		apkFile := apk.
			Directory("/home/builder/packages").
			Directory(arch.ApkArch).
			File(fmt.Sprintf("harbor-cli-%s-r0.apk", m.AppVersion))

		apkFileName := fmt.Sprintf("apk/harbor-cli_%s_%s.apk", m.AppVersion, arch.Arch)
		buildDir = buildDir.WithFile(apkFileName, apkFile)
	}

	return buildDir, nil
}

func apkbuild(ver string) string {
	return fmt.Sprintf(`# APKBUILD
pkgname=harbor-cli
pkgver=%s
pkgrel=0
pkgdesc="Harbor CLI — a command-line interface for interacting with your Harbor container registry."
url="https://github.com/goharbor/harbor-cli"
arch="x86_64 aarch64"
license="Apache-2.0"
depends=""
makedepends=""
source=""
maintainer="Harbor CLI Maintainers <harbor-dev@lists.cncf.io>"
builddir="/build"

package() {
    install -Dm755 "$builddir/harbor-cli" \
        "$pkgdir/usr/bin/harbor-cli"
}
`, ver)
}
