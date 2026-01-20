
![Harbor-CLI Logo_256px](https://github.com/user-attachments/assets/fa18e8f0-a2e4-4462-ab2d-446a88f9edb3)

**Harbor CLI — a command-line interface for interacting with your Harbor container registry. A streamlined, user-friendly alternative to the WebUI, as your daily driver or for scripting and automation.**

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/harbor-cli)](https://artifacthub.io/packages/search?repo=harbor-cli)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fgoharbor%2Fharbor-cli.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fgoharbor%2Fharbor-cli?ref=badge_shield)
[![codecov](https://codecov.io/gh/goharbor/harbor-cli/branch/main/graph/badge.svg)](https://codecov.io/gh/goharbor/harbor-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/goharbor/harbor-cli)](https://goreportcard.com/report/github.com/goharbor/harbor-cli)

# Scope 🧪

1. CLI alternative to the WebUI
2. Tool for scripting and automation of common repeatable Harbor tasks running on your machine or inside your pipeline

# Features
The project's first goal is to reach WebUI parity.

```
✅ project       Manage projects  
✅ repo          Manage repositories  
✅ artifact      Manage artifacts  
✅ label         Manage labels  
✅ tag           Manage tags   
✅ quota         Manage quotas  
✅ webhook       Manage webhook policies 
✅ robot         Robot Account 

✅ login         Log in to Harbor registry  
✅ user          Manage users  

✅ registry      Manage registries
✅ replication   Manage replication

✅ config        Manage the config of the Harbor CLI
✅ cve-allowlist Manage system CVE allowlist
✅ health        Get the health status of Harbor components
✅ preheat       Manage preheat provider instances
✅ info          Display detailed Harbor system, statistics, and CLI environment information

✅ Scanner       scanner commands CRUD
✅ schedule      Schedule jobs in Harbor

❌ Vulnerability Dashboard

❌ Distribution

❌ GC

❌ Job Service   Dashbaord

❌ Auditlog      Auditlogs dashboard

✅ completion    Generate the autocompletion script for the specified shell
✅ help          Help about any command
✅ version       Version of Harbor CLI

```

# Installation

## Container 

Running Harbor CLI as a container is simple. Use the following command to get started:

```shell
docker run -ti --rm -v $HOME/.config/harbor-cli/config.yaml:/root/.config/harbor-cli/config.yaml \
  -e HARBOR_ENCRYPTION_KEY=$(echo "ThisIsAVeryLongPassword" | base64) \
  registry.goharbor.io/harbor-cli/harbor-cli \
  --help
```
Use the `HARBOR_ENCRYPTION_KEY` container environment variable as a base64-encoded 32-byte key for AES-256 encryption. This securely stores your harbor login password.

If you intend
to run the CLI as a container,it is advised
to set the following environment variables and to create an alias
and append the alias to your .zshrc or .bashrc file

```shell
echo "export HARBOR_CLI_CONFIG=\$HOME/.config/harbor-cli/config.yaml" >> ~/.zshrc
echo "export HARBOR_ENCRYPTION_KEY=\$(cat <path_to_32bit_private_key_file> | base64)" >> ~/.zshrc
echo "alias harbor='docker run -ti --rm -v \$HARBOR_CLI_CONFIG:/root/.config/harbor-cli/config.yaml -e HARBOR_ENCRYPTION_KEY=\$HARBOR_ENCRYPTION_KEY registry.goharbor.io/harbor-cli/harbor-cli'" >> ~/.zshrc 
source ~/.zshrc # or restart your terminal
```

## Linux, macOS and Windows

On Linux and macOS, you can use Homebrew:

```bash
brew install harbor-cli
```

Otherwise, you can download the binary from the [releases page](https://github.com/goharbor/harbor-cli/releases).

## Add the Harbor CLI to your Container Image

Using Curl or Wget isn't needed if you want to 
add the Harbor CLI to your container.
Instead, we recommend copying the Harbor CLI from our official image
by using the following Dockerfile:

```Dockerfile
#...
COPY --from=registry.goharbor.io/harbor-cli/harbor-cli:latest /harbor /usr/local/bin/harbor
# --chown and --chmod flags can be used to set the permissions
```



# Example Commands💡
```bash
>./harbor    

Official Harbor CLI

Usage:
  harbor [command]

Examples:

// Base command:
harbor

// Display help about the command:
harbor help


Available Commands:
  artifact      Manage artifacts
  completion    Generate the autocompletion script for the specified shell
  config        Manage the config of the Harbor CLI
  cve-allowlist Manage system CVE allowlist
  health        Get the health status of Harbor components
  help          Help about any command
  info          Show the current credential information
  instance      Manage preheat provider instances in Harbor
  label         Manage labels in Harbor
  login         Log in to Harbor registry
  project       Manage projects and assign resources to them
  registry      Manage registries
  repo          Manage repositories
  schedule      Schedule jobs in Harbor
  tag           Manage tags in Harbor registry
  user          Manage users
  version       Version of Harbor CLI

Flags:
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -h, --help                   help for harbor
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output

Use "harbor [command] --help" for more information about a command.


```

#### Config Management

##### Hierarchy
  Use the `--config` flag to specify a custom configuration file path (the highest priority).
  
```bash
  harbor --config /path/to/custom/config.yaml artifact list
  ```

  If `--config` is not provided, Harbor CLI checks the `HARBOR_CLI_CONFIG` environment variable for the config file path.

  ```bash
  export HARBOR_CLI_CONFIG=/path/to/custom/config.yaml
  harbor artifact list
  ```

  If neither is set, it defaults to `$XDG_CONFIG_HOME/harbor-cli/config.yaml` or `$HOME/.config/harbor-cli/config.yaml` if `XDG_CONFIG_HOME` is unset.
  ```bash
  harbor artifact list
  ```  

##### Data Path
  - Data paths are determined by the `XDG_DATA_HOME` environment variable.
  - If `XDG_DATA_HOME` is not set, it defaults to `$HOME/.local/share/harbor-cli/data.yaml`.
  - The data file always contains the path of the latest config used.

##### Config TL;DR
  - `--config` flag > `HARBOR_CLI_CONFIG` environment variable > default XDG config paths.
  - Environment variables override default settings, and the `--config` flag takes precedence over both environment variables and defaults.
  - The data file always contains the path of the latest config used.


#### Log in to Harbor Registry

```bash
harbor login demo.goharbor.io -u harbor-cli -p Harbor12345
```

#### Create a New Project

```bash
harbor project create
```

#### List all Projects

```bash
harbor project list

# output
┌──────────────────────────────────────────────────────────────────────────────────────────┐
│  Project Name  Access Level  Type          Repo Count    Creation Time                   │
│ ──────────────────────────────────────────────────────────────────────────────────────── │
│  library       public        project       0             1 hour ago                      │
└──────────────────────────────────────────────────────────────────────────────────────────┘
```

#### List all Repository in a Project

```bash
harbor repo list

# output
┌────────────────────────────────────────────────────────────────────────────────────────┐
│  Name                      Artifacts     Pulls         Last Modified Time              │
│ ────────────────────────────────────────────────────────────────────────────────────── │
│  library/harbor-cli        1             0             0 minute ago                    │
└────────────────────────────────────────────────────────────────────────────────────────┘
```

# Supported Platforms

Platform | Status
--|--
Linux | ✅
macOS | ✅
Windows | ✅



# Build From Source

Make sure you have the latest [Dagger](https://docs.dagger.io/) installed in your system. 

#### Using Dagger

```bash
git clone https://github.com/goharbor/harbor-cli.git && cd harbor-cli
dagger call build-dev --platform darwin/arm64 export --path=./harbor-cli
./harbor-dev --help
```

If golang is installed in your system, you can also build the project using the following commands:

```bash
git clone https://github.com/goharbor/harbor-cli.git && cd harbor-cli
go build -o harbor-cli cmd/harbor/main.go
```

# Version Compatibility With Harbor

At the moment, the Harbor CLI is developed and tested with Harbor 2.13.
The CLI should work with versions prior to 2.13,
but not all functionalities may be available or work as expected.

Harbor <2.0.0 is not supported.



# Community

* **Twitter:** [@project_harbor](https://twitter.com/project_harbor)
* **User Group:** Join Harbor user email group: [harbor-users@lists.cncf.io](https://lists.cncf.io/g/harbor-users) to get update of Harbor's news, features, releases, or to provide suggestion and feedback.
* **Developer Group:** Join Harbor developer group: [harbor-dev@lists.cncf.io](https://lists.cncf.io/g/harbor-dev) for discussion on Harbor development and contribution.
* **Slack:** Join Harbor's community for discussion and ask questions: [Cloud Native Computing Foundation](https://slack.cncf.io/), channel: [#harbor](https://cloud-native.slack.com/messages/harbor/), [#harbor-dev](https://cloud-native.slack.com/messages/harbor-dev/) and [#harbor-cli](https://cloud-native.slack.com/archives/C078LCGU9K6).
* **Community Calls:** Every Tuesday at 15:00 CET/CEST or 19:30 IST - [Join Meeting](https://zoom.us/j/99658352431)

# License

This project is licensed under the Apache 2.0 License. See the [LICENSE](https://github.com/goharbor/harbor-cli/blob/main/LICENSE) file for details.


[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fgoharbor%2Fharbor-cli.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fgoharbor%2Fharbor-cli?ref=badge_large)

# Acknowledgements

This project is maintained by the Harbor community. We thank all our contributors and users for their support.

# ❤️ Show your support

For any questions or issues, please open an issue on our [GitHub Issues](https://github.com/goharbor/harbor-cli/issues) page.<br>
Give a ⭐ if this project helped you, Thank YOU!

