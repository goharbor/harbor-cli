# 🛠️ Harbor CLI Dagger Pipeline

We use [Dagger](https://dagger.io) to define a CI/CD pipeline for building, linting, and publishing the [Harbor CLI](https://github.com/goharbor/harbor-cli). 
This README will help beginners understand how to use Dagger in local development and CI workflows.

## Prerequisites

Before you start, ensure you have the following:

1. Dagger: Install the latest version of Dagger. You can check the official documentation for installation steps: [Dagger Installation Guide](https://docs.dagger.io/install).

## Dagger Setup and Development Mode

### Run Dagger Develop

```bash
dagger develop
```

This command will generate the necessary files and configuration for building and running Dagger.


## 📦 Dagger Functions Explained

### 🔧 `BuildDev(platform)`

Builds a development binary for your target platform.

```bash
dagger call build-dev --platform="linux/amd64" export --path=bin/harbor-dev
```

### 🧼 `LintReport()`

Runs `golangci-lint` on your code and saves the report to a file.

```bash
dagger call lint-report export --path=./LintReport.json
```

### 📝 `TestCoverageReport()`

Runs go test coverage tools and creates a report.
```bash
dagger call test-coverage-report export --path=coverage-report.md
```

### ✅ `CheckCoverageThreshold(context, threshold)`

Runs go test coverage tools and creates a report. The total coverage is compared to a threshold that can be set to e.g. 80%.
```bash
dagger call check-coverage-threshold --threshold 80.0 
```

### 🚀 `PublishImage(registry, imageTags)`

Builds and publishes the Harbor CLI image to the given container registry with proper OCI metadata labels.

Before running the command you have to export you registry password

```shell
export REGPASS=Harbor12345
```

```bash
dagger call publish-image \
  --registry=demo.goharbor.io \
  --registry-username=harbor-cli \
  --registry-password=env:REGPASS \
  --imageTags=v0.1.0,latest
```

###  `PublishToWinget(packageId, version, githubToken, installerUrls)`

Automates the submission of Harbor CLI updates to the Windows Package Manager (WinGet) repository. This function uses `wingetcreate` to update the package manifest and automatically submit a pull request to `microsoft/winget-pkgs`.

Before running the command, export your GitHub Personal Access Token (with `public_repo` scope):

```shell
export GITHUB_TOKEN=ghp_yourTokenHere
```

**Basic Usage** (auto-detects installer URLs from GitHub releases):

```bash
dagger call publish-to-winget \
  --package-id="GoHarbor.Harbor" \
  --version="0.0.11" \
  --github-token=env:GITHUB_TOKEN
```

**Advanced Usage** (specify custom installer URLs):

```bash
dagger call publish-to-winget \
  --package-id="GoHarbor.Harbor" \
  --version="0.0.11" \
  --github-token=env:GITHUB_TOKEN \
  --installer-urls="https://github.com/goharbor/harbor-cli/releases/download/v0.0.11/harbor-cli_0.0.11_windows_amd64.zip" \
  --installer-urls="https://github.com/goharbor/harbor-cli/releases/download/v0.0.11/harbor-cli_0.0.11_windows_arm64.zip"
```

This will:
1. Download the latest `wingetcreate` tool
2. Update the WinGet manifest with the new version and installer URLs
3. Automatically submit a PR to the `microsoft/winget-pkgs` repository

**Requirements:**
- GitHub Personal Access Token with `public_repo` scope
- Valid installer URLs (must be publicly accessible)
- Existing package in the WinGet repository
- **Windows Docker host** or Windows container support (for CI/CD, use `runs-on: windows-latest` in GitHub Actions)

**Note:** This function requires Windows containers because `wingetcreate` is a Windows-only tool. On Linux/macOS hosts, this will fail. In GitHub Actions, ensure your workflow uses a Windows runner.

###  `PublishToWingetDryRun(packageId, version, installerUrls)`

Test the WinGet publishing logic without requiring Windows containers. This shows exactly what command would be executed.

```bash
dagger call publish-to-winget-dry-run \
  --package-id="GoHarbor.Harbor" \
  --version="0.0.11"
```

Output example:
```
📦 WinGet Publishing Dry Run
========================================
Package ID: GoHarbor.Harbor
Version: 0.0.11
Installer URLs:
https://github.com/goharbor/harbor-cli/releases/download/v0.0.11/harbor-cli_0.0.11_windows_amd64.zip
https://github.com/goharbor/harbor-cli/releases/download/v0.0.11/harbor-cli_0.0.11_windows_arm64.zip

Command that would be executed:
wingetcreate.exe update GoHarbor.Harbor --version 0.0.11 --urls "..." --submit --token $env:GITHUB_TOKEN
```

---

## ⚙️ Configuration Constants

Dagger uses these constant versions (you can modify them as needed):

```go
const (
  GO_VERSION           = "1.24.2"
  GOLANGCILINT_VERSION = "v2.1.2"
  SYFT_VERSION         = "v1.9.0"
  GORELEASER_VERSION   = "v2.3.2"
)
```

---

## 💡 Tips for Beginners

- Every container step is **reproducible** you can build locally or in GitHub Actions without changes.
- Use Dagger to cache Go builds and lint output, speeding up re-runs.

---

## 📚 References

- [Dagger Go SDK Docs](https://pkg.go.dev/dagger.io/dagger)
- [golangci-lint](https://golangci-lint.run/)
- [Goreleaser](https://goreleaser.com/)
