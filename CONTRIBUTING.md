# Contributing to Harbor CLI

Thank you for your interest in contributing to the Harbor CLI project!
We welcome contributions of all kinds, from bug fixes and documentation improvements to new features and suggestions.

## Overview

The **Harbor CLI** is a powerful command-line tool to interact with the [Harbor container registry](https://goharbor.io/). It's built in Go and helps users manage Harbor resources like projects, registries, artifacts, and more — directly from their terminal.

## Getting Started

### Run using Container

Running Harbor CLI as a container is simple. Use the following command to get started:

```shell
docker run -ti --rm -v $HOME/.config/harbor-cli:/root/.config/harbor-cli \
  -e HARBOR_ENCRYPTION_KEY=$(echo "ThisIsAVeryLongPassword" | base64) \
  registry.goharbor.io/harbor-cli/harbor-cli \
  --help
```
Use the `HARBOR_ENCRYPTION_KEY` container environment variable as a base64-encoded 32-byte key for AES-256 encryption. This securely stores your harbor login password.

If you intend to run the CLI as a container, it is advised
to set the following environment variables and to create an alias
and append the alias to your .zshrc or .bashrc file

```shell
echo "export HARBOR_CLI_CONFIG=\$HOME/.config/harbor-cli" >> ~/.zshrc
echo "export HARBOR_ENCRYPTION_KEY=\$(cat <path_to_32bit_private_key_file> | base64)" >> ~/.zshrc
echo "alias harbor='docker run -ti --rm -v \$HARBOR_CLI_CONFIG:/root/.config/harbor-cli -e HARBOR_ENCRYPTION_KEY=\$HARBOR_ENCRYPTION_KEY registry.goharbor.io/harbor-cli/harbor-cli'" >> ~/.zshrc 
source ~/.zshrc # or restart your terminal
```

### Build from Source

Make sure [Go](https://go.dev/) is installed (≥ v1.24).

```bash
git clone https://github.com/goharbor/harbor-cli.git && cd harbor-cli
go build -o harbor-cli cmd/harbor/main.go
./harbor-cli --help
```

Alternatively, use [Dagger](https://docs.dagger.io/) for isolated builds:

```bash
dagger call build-dev --platform darwin/arm64 export --path=./harbor-cli
./harbor-dev --help
```

## Project Structure

```
..
├── cmd/harbor/           # Entry point (main.go) and all CLI commands (Cobra-based)
├── pkg/                  # Shared utilities and internal packages used across commands
├── doc/                 # Project documentation
├── test/                 # CLI tests and test data
├── .github/              # GitHub workflows and issue templates
├── go.mod / go.sum       # Go module dependencies
└── README.md             # Project overview and usage
```

## How to Contribute

### 1. [Fork](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo) and Clone

```bash
git clone https://github.com/your-username/harbor-cli.git
cd harbor-cli
```

### 2. Create Your Feature Branch

```bash
git checkout -b feat/<your-feature-name>
```

### 3. Make Your Changes

Follow coding and formatting guidelines.

### 4. Test Locally

Ensure your changes work as expected.

```bash
gofmt -s -w .
dagger call build-dev --platform darwin/arm64 export --path=./harbor-cli  #Recommended
./harbor-dev --help
```

If dagger is not installed in your system, you can also build the project using the following commands:

```bash
gofmt -s -w .
go build -o ./bin/harbor-cli cmd/harbor/main.go
./bin/harbor-cli --help
```

### 5. Update Documentation 

Before committing, **always regenerate the documentation** if you've made any code changes or added new commands:

```bash
dagger call run-doc export --path=./doc
```

### 6. Commit with a clear message

```bash
git commit -s -m "feat(project): add delete command for project resources"
```

### 7. Push to your own fork

```bash
git push origin feat/<your-feature-name>
```

### 8. Create a dedicated issue **first**

Before opening a pull request, please [open a dedicated issue](https://docs.github.com/en/issues/tracking-your-work-with-issues/using-issues/creating-an-issue) describing the problem you're solving or the feature you're proposing. This gives maintainers a chance to provide early feedback and helps avoid duplicated effort.

**Pull requests without a linked issue will not be reviewed and may be closed.**

### 9. Open a Pull Request

Once your issue has been acknowledged, [open a pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request) and [link it to your issue](https://docs.github.com/en/issues/tracking-your-work-with-issues/using-issues/linking-a-pull-request-to-an-issue).

A good PR includes:
- A clear description of **what** changed and **why**
- Evidence that you tested your changes (screenshots, terminal output, etc.)
- Well-structured, readable commits

> **A note on AI-generated contributions:** We appreciate the intent, but please don't submit PRs that are purely AI-generated without your own understanding and review. Drive-by PRs — especially bulk or low-effort ones produced by AI tools — add review burden for maintainers and will be closed. If you use AI to assist your work, that's fine — just make sure you understand every line you're submitting and can speak to the changes in review.

## 🧪 Running Tests

> ✅ Note: Add your CLI or unit tests to the `test/` directory.

```bash
go test ./...
```

## 🧹 Code Guidelines

- Use `go fmt ./...` to format your code.
- Use descriptive commit messages:
  - `feat`: New feature
  - `fix`: Bug fix
  - `docs`: Documentation only
  - `test`: Adding or updating tests
  - `refactor`: Code cleanup
  - `chore`: Maintenance tasks

## 📬 Communication

- **Slack:** Join us in [#harbor-cli](https://cloud-native.slack.com/messages/harbor-cli/)
- **Issues:** Use [GitHub Issues](https://github.com/goharbor/harbor-cli/issues) for bugs, ideas, or questions.
- **Mailing List:**
  - Users: [harbor-users@lists.cncf.io](https://lists.cncf.io/g/harbor-users)
  - Devs: [harbor-dev@lists.cncf.io](https://lists.cncf.io/g/harbor-dev)

## 📄 License

All contributions are under the [Apache 2.0 License](./LICENSE).

---

**Thank you for contributing to Harbor CLI! Your work helps improve the Harbor ecosystem for everyone. 🙌**
