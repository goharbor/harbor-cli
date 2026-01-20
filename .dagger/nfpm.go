package main

import (
	"context"
	"fmt"

	"dagger/harbor-cli/internal/dagger"
)

func (m *HarborCli) NfpmBuild(ctx context.Context,
	buildDir *dagger.Directory,
	// +ignore=[".gitignore"]
	// +defaultPath="."
	// +optional
	source *dagger.Directory,
) (*dagger.Directory, error) {
	if !m.IsInitialized {
		err := m.init(ctx, source)
		if err != nil {
			return nil, err
		}
	}

	archs := []string{"amd64", "arm64"}
	pkgs := []string{"deb", "rpm"}

	for _, pkg := range pkgs {
		out := dag.Directory()
		for _, arch := range archs {
			fileName := fmt.Sprintf("harbor-cli_%s_%s.%s", m.AppVersion, arch, pkg)

			out = TemplatedYML(out, arch, m.AppVersion, fmt.Sprintf("harbor-cli_%s_%s_%s", m.AppVersion, "linux", arch))

			pkgFile := dag.Container().
				From("goreleaser/nfpm").
				WithMountedFile("/nfpm.yml", out.File("/nfpm.yml")).
				WithMountedDirectory("/input", buildDir).
				WithWorkdir("/input").
				WithExec([]string{
					"nfpm",
					"pkg",
					"--config", "/nfpm.yml",
					"--packager", pkg,
					"--target", fmt.Sprintf("/%s", fileName),
				}).
				File(fmt.Sprintf("/%s", fileName))

			out = out.WithFile(fileName, pkgFile)
		}

		buildDir = buildDir.WithDirectory(fmt.Sprintf("/%s", pkg), out)
	}

	return buildDir, nil
}

func TemplatedYML(out *dagger.Directory, arch string, appV string, filename string) *dagger.Directory {
	out = out.WithNewFile("nfpm.yml", fmt.Sprintf(`
name: harbor-cli
arch: %s
platform: linux
version: %s
section: default
priority: extra
maintainer: "Harbor CLI Maintainers <harbor-dev@lists.cncf.io>"
description: "Harbor CLI — a command-line interface for interacting with your Harbor container registry."
license: Apache 2.0 
contents:
  - src: ./bin/%s
    dst: /usr/local/bin/harbor-cli
`, arch, appV, filename))

	return out
}
