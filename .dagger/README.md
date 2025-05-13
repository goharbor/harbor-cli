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
dagger call test-coverage-report export --path=CoverageReport.out
```

### ✅ `CheckCoverageThreshold()`

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
