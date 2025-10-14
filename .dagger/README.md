# 🛠️ Harbor CLI — Dagger Pipeline

We use [Dagger](https://dagger.io) to define a **modular and reproducible CI/CD pipeline** for building, linting, testing, and publishing the [Harbor CLI](https://github.com/goharbor/harbor-cli).  
This README provides a clear reference for contributors and maintainers to understand, run, and extend the pipeline locally or in CI.

---

## 🚧 Prerequisites

Before using the pipeline, make sure you have:

1. **Dagger CLI** — Install the latest version from the official docs:  
   👉 [Dagger Installation Guide](https://docs.dagger.io/install)
2. **Go** — Installed according to the version specified in the project’s `go.mod`.
3. **Docker** — Required if you’re publishing images.

---

## ⚙️ Setup and Development Mode

### Run Dagger in Development Mode

To start the Dagger session and enable live code reloads:

```bash
dagger develop
```

This command prepares the environment for pipeline development and local testing.

## 📦 Dagger Functions Overview

| **Name**                       | **Description**                                                                                 |
|--------------------------------|-------------------------------------------------------------------------------------------------|
| `lint`                         | Runs `golangci-lint` and prints the report as a string to stdout.                              |
| `lint-report`                  | Runs `golangci-lint` and writes the lint report to a file.                                     |
| `pipeline`                     | Executes the **full CI/CD pipeline** including build, test, lint, and publish stages.          |
| `run-doc`                      | Generates CLI documentation and returns the directory containing generated files.              |
| `test`                         | Runs all Go tests in the repository.                                                           |
| `test-report`                  | Executes Go tests and outputs a structured JSON test report.                                   |
| `test-coverage`                | Runs Go tests with coverage tracking.                                                          |
| `test-coverage-report`         | Processes coverage data and returns a formatted Markdown report.                               |
| `vulnerability-check`          | Runs `govulncheck` to detect known vulnerabilities in dependencies.                            |
| `vulnerability-check-report`   | Runs `govulncheck` and saves results to a file (`vulnerability-check.report`).                  |
| `build-dev`                    | Create build of Harbor CLI for local testing and development|

---

## 🧩 Example Usage

Below are some common commands to run specific Dagger functions locally:

```bash
# Development build for binaries

dagger call build-dev --source=. --platform="linux/amd64" export --path=bin/harbor-dev

# Print report to stdout
dagger call lint

# Save report to a file
dagger call lint-report export --path=LintReport.json

# Run Tests
dagger call test

# Generate a JSON Report
dagger call test-report export --path=TestReport.json

# Test Coverage
dagger call test-coverage

# Generate a Markdown Report
dagger call test-coverage-report export --path=coverage-report.md

# Vulnerability Check 
dagger call vulnerability-check 

# Generate a Report
dagger call vulnerability-check-report export --path=vuln.report

# Generate CLI docs 
dagger call run-doc export --path=docs/cli 
```


## 💡 Tips for Contributors

- Every step in Dagger is **deterministic and reproducible** — what you run locally is identical to CI.
- Use Dagger’s built-in caching to accelerate Go builds, lint runs, and dependency installs.
- Modular functions let you run only what you need, improving iteration speed and debugging efficiency.
- Prefer using `dagger develop` for fast iteration and testing new steps before committing.
- Store output reports (lint, test, coverage) under a consistent `/reports` directory for easier CI integration.
- If you modify or add new pipeline steps, document them under **Dagger Functions Overview** to maintain clarity.
- Always validate pipelines with `dagger call pipeline` locally before merging into main.

---

## 📚 References

- [Dagger Go SDK](https://pkg.go.dev/dagger.io/dagger)
- [golangci-lint Docs](https://golangci-lint.run/)
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
