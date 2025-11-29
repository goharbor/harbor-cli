package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dagger/harbor-cli/internal/dagger"
)

func (m *HarborCli) PublishImageAndSign(
	ctx context.Context,
	source *dagger.Directory,
	registry string,
	registryUsername string,
	registryPassword *dagger.Secret,
	imageTags []string,
	// +optional
	githubToken *dagger.Secret,
	// +optional
	actionsIdTokenRequestToken *dagger.Secret,
	// +optional
	actionsIdTokenRequestUrl string,
) (string, error) {
	m.init(ctx, source)

	imageAddrs := m.PublishImage(ctx, registry, registryUsername, imageTags, registryPassword)

	_, err := m.Sign(
		ctx,
		githubToken,
		actionsIdTokenRequestUrl,
		actionsIdTokenRequestToken,
		registryUsername,
		registryPassword,
		imageAddrs[0],
	)
	if err != nil {
		return "", fmt.Errorf("failed to sign image: %w", err)
	}

	fmt.Printf("Signed image: %s\n", imageAddrs)
	return imageAddrs[0], nil
}

func (m *HarborCli) PublishImage(
	ctx context.Context,
	registry, registryUsername string,
	// +optional
	// +default=["latest"]
	imageTags []string,
	registryPassword *dagger.Secret,
) []string {
	version := getVersion(imageTags)
	builders := m.build(ctx, version)
	releaseImages := []*dagger.Container{}

	for i, tag := range imageTags {
		imageTags[i] = strings.TrimSpace(tag)
		if strings.HasPrefix(imageTags[i], "v") {
			imageTags[i] = strings.TrimPrefix(imageTags[i], "v")
		}
	}
	fmt.Printf("provided tags: %s\n", imageTags)

	// Get current time for image creation timestamp
	creationTime := time.Now().UTC().Format(time.RFC3339)

	for _, builder := range builders {
		os, _ := builder.EnvVariable(ctx, "GOOS")
		arch, _ := builder.EnvVariable(ctx, "GOARCH")

		if os != "linux" {
			continue
		}

		ctr := dag.Container(dagger.ContainerOpts{Platform: dagger.Platform(os + "/" + arch)}).
			From("alpine:latest").
			WithWorkdir("/").
			WithFile("/harbor", builder.File("./harbor")).
			WithExec([]string{"ls", "-al"}).
			WithExec([]string{"./harbor", "version"}).
			// Add required metadata labels for ArtifactHub
			WithLabel("org.opencontainers.image.created", creationTime).
			WithLabel("org.opencontainers.image.description", "Harbor CLI - A command-line interface for CNCF Harbor, the cloud native registry!").
			WithLabel("io.artifacthub.package.readme-url", "https://raw.githubusercontent.com/goharbor/harbor-cli/main/README.md").
			WithLabel("org.opencontainers.image.source", "https://github.com/goharbor/harbor-cli").
			WithLabel("org.opencontainers.image.version", version).
			WithLabel("io.artifacthub.package.license", "Apache-2.0").
			WithEntrypoint([]string{"/harbor"})
		releaseImages = append(releaseImages, ctr)
	}

	imageAddrs := []string{}
	for _, imageTag := range imageTags {
		addr, err := dag.Container().WithRegistryAuth(registry, registryUsername, registryPassword).
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

func (m *HarborCli) build(
	ctx context.Context,
	version string,
) []*dagger.Container {
	var builds []*dagger.Container

	fmt.Println("🛠️  Building with Dagger...")
	oses := []string{"linux", "darwin", "windows"}
	arches := []string{"amd64", "arm64"}

	// temp container with git installed
	temp := dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/src", m.Source).
		// --no-cache option is to avoid caching the apk package index
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithWorkdir("/src")

	gitCommit, _ := temp.WithExec([]string{"git", "rev-parse", "--short", "HEAD", "--always"}).Stdout(ctx)
	buildTime := time.Now().UTC().Format(time.RFC3339)
	ldflagsArgs := fmt.Sprintf(`-X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.Version=%s
						  -X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.GoVersion=%s
						  -X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.BuildTime=%s
						  -X github.com/goharbor/harbor-cli/cmd/harbor/internal/version.GitCommit=%s
				`, version, m.GoVersion, buildTime, gitCommit)

	for _, goos := range oses {
		for _, goarch := range arches {
			bin_path := fmt.Sprintf("build/%s/%s/", goos, goarch)
			builder := dag.Container().
				From("golang:"+m.GoVersion+"-alpine").
				WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+m.GoVersion)).
				WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
				WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+m.GoVersion)).
				WithEnvVariable("GOCACHE", "/go/build-cache").
				WithMountedDirectory("/src", m.Source).
				WithWorkdir("/src").
				WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch).
				WithExec([]string{"go", "build", "-ldflags", ldflagsArgs, "-o", bin_path + "harbor", "/src/cmd/harbor/main.go"}).
				WithWorkdir(bin_path).
				WithExec([]string{"ls"}).
				WithEntrypoint([]string{"./harbor"})

			builds = append(builds, builder)
		}
	}
	return builds
}

// Sign signs a container image using Cosign, works also with GitHub Actions
func (m *HarborCli) Sign(ctx context.Context,
	// +optional
	githubToken *dagger.Secret,
	// +optional
	actionsIdTokenRequestUrl string,
	// +optional
	actionsIdTokenRequestToken *dagger.Secret,
	registryUsername string,
	registryPassword *dagger.Secret,
	imageAddr string,
) (string, error) {
	registryPasswordPlain, _ := registryPassword.Plaintext(ctx)

	cosing_ctr := dag.Container().From("cgr.dev/chainguard/cosign")

	// If githubToken is provided, use it to sign the image
	if githubToken != nil {
		if actionsIdTokenRequestUrl == "" || actionsIdTokenRequestToken == nil {
			return "", fmt.Errorf("actionsIdTokenRequestUrl (exist=%s) and actionsIdTokenRequestToken (exist=%t) must be provided when githubToken is provided", actionsIdTokenRequestUrl, actionsIdTokenRequestToken != nil)
		}
		fmt.Printf("Setting the ENV Vars GITHUB_TOKEN, ACTIONS_ID_TOKEN_REQUEST_URL, ACTIONS_ID_TOKEN_REQUEST_TOKEN to sign with GitHub Token")
		cosing_ctr = cosing_ctr.WithSecretVariable("GITHUB_TOKEN", githubToken).
			WithEnvVariable("ACTIONS_ID_TOKEN_REQUEST_URL", actionsIdTokenRequestUrl).
			WithSecretVariable("ACTIONS_ID_TOKEN_REQUEST_TOKEN", actionsIdTokenRequestToken)
	}

	return cosing_ctr.WithSecretVariable("REGISTRY_PASSWORD", registryPassword).
		WithExec([]string{"cosign", "env"}).
		WithExec([]string{
			"cosign", "sign", "--yes", "--recursive",
			"--registry-username", registryUsername,
			"--registry-password", registryPasswordPlain,
			imageAddr,
			"--timeout", "1m",
		}).Stdout(ctx)
}

func getVersion(tags []string) string {
	for _, tag := range tags {
		if strings.HasPrefix(tag, "v") {
			return tag
		}
	}
	return "latest"
}
