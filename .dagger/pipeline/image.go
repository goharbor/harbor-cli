package pipeline

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dagger/harbor-cli/internal/dagger"
)

// PublishImage publishes a container image to a registry with a specific tag and signs it using Cosign.
func (s *Pipeline) PublishImage(
	ctx context.Context,
	dist *dagger.Directory,
	registry, registryUsername string,
	// +optional
	// +default=["latest"]
	imageTags []string,
	registryPassword *dagger.Secret,
) []string {
	version := s.appVersion
	archs := []string{"amd64", "arm64"}
	releaseImages := []*dagger.Container{}

	for i, tag := range imageTags {
		imageTags[i] = strings.TrimSpace(tag)
		imageTags[i] = strings.TrimPrefix(imageTags[i], "v")
	}
	fmt.Printf("provided tags: %s\n", imageTags)

	// Get current time for image creation timestamp
	creationTime := time.Now().UTC().Format(time.RFC3339)

	for _, arch := range archs {
		binName := fmt.Sprintf("harbor-cli_%s_linux_%s", s.appVersion, arch)

		ctr := s.dag.Container(dagger.ContainerOpts{Platform: dagger.Platform("linux/" + arch)}).
			From("alpine:latest").
			WithWorkdir("/src").
			WithFile("./harbor", dist.File(fmt.Sprintf("linux/%s", binName))).
			WithExec([]string{"ls", "-al"}).
			WithExec([]string{"./harbor", "version"}).
			// Add required metadata labels for ArtifactHub
			WithLabel("org.opencontainers.image.created", creationTime).
			WithLabel("org.opencontainers.image.description", "Harbor CLI - A command-line interface for CNCF Harbor, the cloud native registry!").
			WithLabel("io.artifacthub.package.readme-url", "https://raw.githubusercontent.com/goharbor/harbor-cli/main/README.md").
			WithLabel("org.opencontainers.image.source", "https://github.com/goharbor/harbor-cli").
			WithLabel("org.opencontainers.image.version", version).
			WithLabel("io.artifacthub.package.license", "Apache-2.0").
			WithEntrypoint([]string{"./harbor"})
		releaseImages = append(releaseImages, ctr)
	}

	imageAddrs := []string{}
	for _, imageTag := range imageTags {
		addr, err := s.dag.Container().WithRegistryAuth(registry, registryUsername, registryPassword).
			Publish(ctx,
				fmt.Sprintf("%s/%s/harbor-cli:%s", registry, registryUsername, imageTag),
				dagger.ContainerPublishOpts{PlatformVariants: releaseImages},
			)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Published image address: %s\n", addr)
		imageAddrs = append(imageAddrs, addr)
	}
	return imageAddrs
}
