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
	// +optional
	buildDir *dagger.Directory,
	// +ignore=[".gitignore"]
	// +defaultPath="."
	source *dagger.Directory,
	registry string,
	registryUsername string,
	registryPassword *dagger.Secret,
	imageTags string,
	// +optional
	githubToken *dagger.Secret,
	// +optional
	actionsIdTokenRequestToken *dagger.Secret,
	// +optional
	actionsIdTokenRequestUrl *dagger.Secret,
) (string, error) {
	if !m.IsInitialized {
		err := m.init(ctx, source)
		if err != nil {
			return "", err
		}
	}

	imageAddrs, err := m.PublishImage(ctx, registry, registryUsername, strings.Split(imageTags, ","), buildDir, source, registryPassword)
	if err != nil {
		return "", err
	}

	for _, addr := range imageAddrs {
		// Generate SBOM (SPDX JSON) for the published image
		sbom := m.GenerateSBOM(
			ctx,
			addr, 
			registryUsername, 
			registryPassword,
		);

		// Attest SBOM to the image using Cosign in-toto attestation
		if _, err := m.AttestSBOM(
			ctx,
			sbom,
			addr,
			registryUsername,
			registryPassword,
			githubToken,
			actionsIdTokenRequestUrl,
			actionsIdTokenRequestToken,
		); err != nil {
			return "", fmt.Errorf("failed to attest SBOM for image %s: %w", addr, err)
		}
		fmt.Printf("Attested SBOM for image: %s\n", addr)

		_, err = m.Sign(
			ctx,
			githubToken,
			actionsIdTokenRequestUrl,
			actionsIdTokenRequestToken,
			registryUsername,
			registryPassword,
			addr,
		)
		if err != nil {
			return "", fmt.Errorf("failed to sign image %s: %w", addr, err)
		}
		fmt.Printf("Signed image: %s\n", addr)
	}

	return imageAddrs[0], nil
}

func (m *HarborCli) PublishImage(
	ctx context.Context,
	registry, registryUsername string,
	// +optional
	// +default=["latest"]
	imageTags []string,
	// +optional
	buildDir *dagger.Directory,
	// +ignore=[".gitignore"]
	// +defaultPath="."
	source *dagger.Directory,
	registryPassword *dagger.Secret,
) ([]string, error) {
	if !m.IsInitialized {
		err := m.init(ctx, source)
		if err != nil {
			return []string{}, err
		}
	}

	version := getVersion(imageTags)
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

	// If the buildDir is not provided, build new binaries ones
	if buildDir == nil {
		buildDir = dag.Directory()

		builders := m.build(ctx, version)

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
	} else { // If buildDir is provided, use existing binaries
		archs := []string{"amd64", "arm64"}

		for _, arch := range archs {
			filepath := fmt.Sprintf("bin/harbor-cli_%s_linux_%s", m.AppVersion, arch)

			ctr := dag.Container(dagger.ContainerOpts{Platform: dagger.Platform("linux/" + arch)}).
				From("alpine:latest").
				WithWorkdir("/").
				WithFile("/harbor", buildDir.File(filepath)).
				WithExec([]string{"ls", "-al"}).
				WithExec([]string{"chmod", "+x", "/harbor"}).
				WithExec([]string{"uname", "-m"}).
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
	}

	imageAddrs := []string{}
	for _, imageTag := range imageTags {
		addr, err := dag.Container().WithRegistryAuth(registry, registryUsername, registryPassword).
			Publish(ctx,
				fmt.Sprintf("%s/%s/harbor-cli:%s", registry, "harbor-cli", imageTag),
				dagger.ContainerPublishOpts{PlatformVariants: releaseImages},
			)
		if err != nil {
			return []string{}, err
		}

		fmt.Printf("Published image address: %s\n", addr)
		imageAddrs = append(imageAddrs, addr)
	}

	return imageAddrs, nil
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
	actionsIdTokenRequestUrl *dagger.Secret,
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
		if actionsIdTokenRequestUrl == nil || actionsIdTokenRequestToken == nil {
			return "", fmt.Errorf("actionsIdTokenRequestUrl (exist=%s) and actionsIdTokenRequestToken (exist=%t) must be provided when githubToken is provided", actionsIdTokenRequestUrl, actionsIdTokenRequestToken != nil)
		}
		fmt.Printf("Setting the ENV Vars GITHUB_TOKEN, ACTIONS_ID_TOKEN_REQUEST_URL, ACTIONS_ID_TOKEN_REQUEST_TOKEN to sign with GitHub Token")
		cosing_ctr = cosing_ctr.WithSecretVariable("GITHUB_TOKEN", githubToken).
			WithSecretVariable("ACTIONS_ID_TOKEN_REQUEST_URL", actionsIdTokenRequestUrl).
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

// GenerateSBOM uses Syft to create an SPDX JSON SBOM for a given image and returns it as a file
func (m *HarborCli) GenerateSBOM(
	ctx context.Context,
	imageAddr string,
	registryUsername string,
	registryPassword *dagger.Secret,
) *dagger.File {
	syftCtr := dag.Container().
		From("anchore/syft:latest").
		WithSecretVariable("SYFT_REGISTRY_AUTH_PASSWORD", registryPassword).
		WithEnvVariable("SYFT_REGISTRY_AUTH_USERNAME", registryUsername).
		// Output SPDX JSON to a known path
		WithExec([]string{"syft", imageAddr, "-o", "spdx-json=/sbom.spdx.json"})

	return syftCtr.File("/sbom.spdx.json")
}

// AttestSBOM attaches an in-toto attestation (SBOM predicate) to the image using Cosign
func (m *HarborCli) AttestSBOM(
	ctx context.Context,
	sbomFile *dagger.File,
	imageAddr string,
	registryUsername string,
	registryPassword *dagger.Secret,
	// +optional
	githubToken *dagger.Secret,
	// +optional
	actionsIdTokenRequestUrl *dagger.Secret,
	// +optional
	actionsIdTokenRequestToken *dagger.Secret,
) (string, error) {
	cosignCtr := dag.Container().
		From("cgr.dev/chainguard/cosign").
		WithMountedFile("/sbom.spdx.json", sbomFile).
		WithSecretVariable("REGISTRY_PASSWORD", registryPassword)

	// If githubToken is provided, configure OIDC for keyless signing
	if githubToken != nil {
		if actionsIdTokenRequestUrl == nil || actionsIdTokenRequestToken == nil {
			return "", fmt.Errorf("actionsIdTokenRequestUrl (exist=%v) and actionsIdTokenRequestToken (exist=%t) must be provided when githubToken is provided", actionsIdTokenRequestUrl != nil, actionsIdTokenRequestToken != nil)
		}
		cosignCtr = cosignCtr.
			WithSecretVariable("GITHUB_TOKEN", githubToken).
			WithSecretVariable("ACTIONS_ID_TOKEN_REQUEST_URL", actionsIdTokenRequestUrl).
			WithSecretVariable("ACTIONS_ID_TOKEN_REQUEST_TOKEN", actionsIdTokenRequestToken)
	}

	// Use cosign attest to create in-toto attestation with SPDX JSON predicate
	return cosignCtr.WithExec([]string{
		"cosign", "attest",
		"--yes",
		"--type", "spdxjson",
		"--predicate", "/sbom.spdx.json",
		"--registry-username", registryUsername,
		"--registry-password", registryPassword,
		imageAddr,
		"--timeout", "1m",
	}).Stdout(ctx)
}
