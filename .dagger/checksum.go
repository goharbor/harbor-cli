package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/harbor-cli/internal/dagger"
)

func (m *HarborCli) Checksum(ctx context.Context,
	buildDir *dagger.Directory,
	// +ignore=[".gitignore"]
	// +defaultPath="."
	// +optional
	source *dagger.Directory,
) (*dagger.Directory, error) {
	sums := map[string]string{}
	bins, err := DistBinaries(ctx, dag, buildDir)
	if err != nil {
		return nil, err
	}

	shasum := dag.Container().
		From("alpine").
		WithMountedDirectory("/dist", buildDir).
		WithWorkdir("/dist")

	for _, v := range bins {
		// We Ignore the filepath provided, since it uses the directory structure, ie,
		// archive/bin.tar.gz or rpm/harbor-cli.rpm
		// And Instead when later merging I will strip the prefix
		out, err := shasum.WithExec([]string{"sh", "-c", fmt.Sprintf("sha256sum %s | awk '{print $1}'", v)}).Stdout(ctx)
		if err != nil {
			return nil, err
		}

		split := strings.Split(v, "/")
		sums[out] = split[len(split)-1] // Taking only filename
	}

	content := ""
	for k, v := range sums {
		content += fmt.Sprintf("%s %s\n", strings.TrimSuffix(k, "\n"), v)
	}

	buildDir = buildDir.WithFile("checksums.txt", dag.File("checksums.txt", content))
	return buildDir, err
}
