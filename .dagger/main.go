package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"dagger/harbor-cli/internal/dagger"
)

type HarborCli struct {
	Source        *dagger.Directory // Source Directory where code resides
	AppVersion    string            // Current Version of the app, acquired from git tags
	GoVersion     string            // Go Version used in the current release, acquired from the go.mod file
	IsInitialized bool
}

func (m *HarborCli) init(ctx context.Context, source *dagger.Directory) error {
	out, err := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"git", "describe", "--tags", "--abbrev=0"}).
		Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	goVersion, err := source.File("go.mod").Contents(ctx)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`(?m)^go (\d+\.\d+(\.\d+)?)`)
	match := re.FindStringSubmatch(goVersion)
	if len(match) > 1 {
		m.GoVersion = match[1]
	}

	m.Source = source
	m.AppVersion = strings.TrimSpace(strings.TrimLeft(out, "v"))
	m.IsInitialized = true

	return nil
}

// PublishToScoopDryRun shows what would be published to Scoop without making changes
// This is useful for testing and validation
func (m *HarborCli) PublishToScoopDryRun(
	ctx context.Context,
	version string,
) string {
	baseUrl := fmt.Sprintf("https://github.com/goharbor/harbor-cli/releases/download/v%s", version)

	output := fmt.Sprintf(`Scoop Publishing Dry Run
========================================
Version: %s
Manifest: scoop/harbor-cli.json

Installer URLs:
- 64bit: %s/harbor-cli_%s_windows_amd64.zip
- arm64: %s/harbor-cli_%s_windows_arm64.zip

This will:
1. Download Windows release assets
2. Compute SHA256 hashes
3. Update scoop/harbor-cli.json with new version and hashes
4. Commit and push changes to goharbor/harbor-cli

Note: This is a dry run. Use publish-to-scoop to actually publish.
`, version, baseUrl, version, baseUrl, version)

	return output
}

// PublishToScoop updates the Scoop manifest in the harbor-cli repo
// This uses the same repo approach - no external secrets needed, just GITHUB_TOKEN
func (m *HarborCli) PublishToScoop(
	ctx context.Context,
	version string,
	githubToken *dagger.Secret,
) (string, error) {
	fmt.Println("Publishing to Scoop...")

	// Construct URLs
	baseUrl := fmt.Sprintf("https://github.com/goharbor/harbor-cli/releases/download/v%s", version)
	amd64Url := fmt.Sprintf("%s/harbor-cli_%s_windows_amd64.zip", baseUrl, version)
	arm64Url := fmt.Sprintf("%s/harbor-cli_%s_windows_arm64.zip", baseUrl, version)

	script := fmt.Sprintf(`
set -e

# Download files and compute hashes
echo "Downloading and computing hashes..."
AMD64_HASH=$(curl -sL "%s" | sha256sum | cut -d' ' -f1)
ARM64_HASH=$(curl -sL "%s" | sha256sum | cut -d' ' -f1)

echo "AMD64 Hash: $AMD64_HASH"
echo "ARM64 Hash: $ARM64_HASH"

# Clone the repo
echo "Cloning repository..."
git clone https://x-access-token:${GITHUB_TOKEN}@github.com/goharbor/harbor-cli.git repo
cd repo

# Create a new branch for the update
BRANCH="scoop/update-v%s"
git checkout -b "$BRANCH"

# Update manifest using jq
echo "Updating manifest..."
jq --arg ver "%s" \
   --arg amd64url "%s" \
   --arg arm64url "%s" \
   --arg amd64hash "$AMD64_HASH" \
   --arg arm64hash "$ARM64_HASH" \
   '.version = $ver |
    .architecture."64bit".url = $amd64url |
    .architecture."64bit".hash = $amd64hash |
    .architecture.arm64.url = $arm64url |
    .architecture.arm64.hash = $arm64hash' \
   scoop/harbor-cli.json > scoop/harbor-cli.json.tmp
mv scoop/harbor-cli.json.tmp scoop/harbor-cli.json

# Show updated manifest
echo "Updated manifest:"
cat scoop/harbor-cli.json

# Commit and push the branch
git config user.name "github-actions[bot]"
git config user.email "github-actions[bot]@users.noreply.github.com"
git add scoop/harbor-cli.json
git commit -m "scoop: update harbor-cli to v%s" || { echo "No changes to commit"; exit 0; }
git push origin "$BRANCH"

# Create a PR using GitHub API
echo "Creating pull request..."
curl -s -X POST \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  -H "Accept: application/vnd.github+json" \
  https://api.github.com/repos/goharbor/harbor-cli/pulls \
  -d "{
    \"title\": \"scoop: update harbor-cli to v%s\",
    \"body\": \"Automated Scoop manifest update for v%s.\",
    \"head\": \"$BRANCH\",
    \"base\": \"main\"
  }"

echo "Pull request created successfully!"
`, amd64Url, arm64Url, version, version, amd64Url, arm64Url, version, version, version)

	output, err := dag.Container().
		From("alpine:latest").
		WithSecretVariable("GITHUB_TOKEN", githubToken).
		WithExec([]string{"apk", "add", "--no-cache", "git", "curl", "bash", "jq"}).
		WithExec([]string{"bash", "-c", script}).
		Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("failed to publish to Scoop: %v", err)
	}

	fmt.Println("Scoop manifest updated successfully")
	return output, nil
}
