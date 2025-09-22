package pipeline

import (
	"context"
	"fmt"

	"dagger/harbor-cli/internal/dagger"
)

func (s *Pipeline) NFPMBuild(ctx context.Context, dist *dagger.Directory) (*dagger.Directory, error) {
	archs := []string{"amd64", "arm64"}
	pkgs := []string{"deb", "rpm"}

	for _, pkg := range pkgs {
		out := s.dag.Directory()
		for _, arch := range archs {
			fileName := fmt.Sprintf("harbor-cli_%s_%s.%s", s.appVersion, arch, pkg)

			out = TemplatedYML(out, arch, s.appVersion, fmt.Sprintf("harbor-cli_%s_%s_%s", s.appVersion, "linux", arch))

			pkgFile := s.dag.Container().
				From("goreleaser/nfpm").
				WithMountedFile("/nfpm.yml", out.File("/nfpm.yml")).
				WithMountedDirectory("/input", dist).
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

		dist = dist.WithDirectory(fmt.Sprintf("/%s", pkg), out)
	}

	return dist, nil
}

func TemplatedYML(out *dagger.Directory, arch string, appV string, filename string) *dagger.Directory {
	out = out.WithNewFile("nfpm.yml", fmt.Sprintf(`
name: harbor-cli
arch: %s
platform: linux
version: %s
section: default
priority: extra
maintainer: "NucleoFusion <lakshit.singh.mail@gmail.com>"
description: "Harbor CLI — a command-line interface for interacting with your Harbor container registry."
license: Apache 2.0 
contents:
  - src: ./linux/%s
    dst: /usr/local/bin/harbor-cl
`, arch, appV, filename))

	return out
}
