
![harbor-3](https://github.com/goharbor/harbor-cli/assets/70086051/835ab686-1cce-4ac7-bc57-05a35c2b73cc)

**Welcome to the Harbor CLI project! This powerful command-line tool facilitates seamless interaction with the Harbor container registry. It simplifies various tasks such as creating, updating, and managing projects, registries, and other resources in Harbor.**

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/harbor-cli)](https://artifacthub.io/packages/search?repo=harbor-cli)

# Project Scope 🧪

The Harbor CLI is designed to enhance your interaction with the Harbor container registry. Built on Golang, it offers a user-friendly interface to perform various tasks related to projects, registries, and more. Whether you're creating, updating, or managing resources, the Harbor CLI streamlines your workflow efficiently.

# Project Features 🤯

 🔹 Get details about projects, registries, repositories and more <br>
 🔹 Create new projects, registries, and other resources <br>
 🔹 Delete projects, registries, and other resources <br>
 🔹 Run commands with various flags for enhanced functionality <br>
 🔹 More features coming soon... 🚧


# Installation

## Container 

It is straightforward to use the Harbor CLI as a container. You can run the following command to use the Harbor CLI as a container:

```shell
docker run -ti --rm -v $HOME/.harbor/config.yaml:/root/.harbor/config.yaml registry.goharbor.io/harbor-cli/harbor-cli --help

```

# Add the following command to create an alias and append the alias to your .zshrc or .bashrc file
```shell
echo "alias harbor='docker run -ti --rm -v \$HOME/.harbor/config.yaml:/root/.harbor/config.yaml registry.goharbor.io/harbor-cli/harbor-cli'" >> ~/.zshrc
source ~/.zshrc # or restart your terminal
```


## Linux, MacOS and Windows

Harbor CLI will soon be published on Homebrew.
Meantime, we recommend using Harbor in the Container
or download the binary from the [releases page](https://github.com/goharbor/harbor-cli/releases)



## Add the Harbor CLI to your Container Image

Using Curl or Wget isn't recommended
for adding the Harbor CLI to your container.
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
  artifact    Manage artifacts
  completion  Generate the autocompletion script for the specified shell
  health      Get the health status of Harbor components
  help        Help about any command
  login       Log in to Harbor registry
  project     Manage projects and assign resources to them
  registry    Manage registries
  repo        Manage repositories
  user        Manage users
  version     Version of Harbor CLI

Flags:
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -h, --help                   help for harbor
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output

Use "harbor [command] --help" for more information about a command.


```

#### Config Management

##### Hierachy
  Use the `--config` flag to specify a custom configuration file path (highest priority).
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

Make sure you have latest [Dagger](https://docs.dagger.io/) installed in your system. 

#### Using Dagger
```bash
git clone https://github.com/goharbor/harbor-cli.git && cd harbor-cli
dagger call build-dev --platform darwin/arm64 export --path=./harbor-cli
./harbor-dev --help
```

If golang is installed in your system, you can also build the project using the following commands:

```bash
git clone https://github.com/goharbor/harbor-cli.git
go build -o harbor-cli cmd/harbor/main.go
```

# Version Compatibility With Harbor

At the moment, the Harbor CLI is developed and tested with Harbor 2.11.
The CLI should work with versions prior to 2.11,
but not all functionalities may be available or work as expected.

Harbor <2.0.0 is not supported.



# Community

* **Twitter:** [@project_harbor](https://twitter.com/project_harbor)
* **User Group:** Join Harbor user email group: [harbor-users@lists.cncf.io](https://lists.cncf.io/g/harbor-users) to get update of Harbor's news, features, releases, or to provide suggestion and feedback.
* **Developer Group:** Join Harbor developer group: [harbor-dev@lists.cncf.io](https://lists.cncf.io/g/harbor-dev) for discussion on Harbor development and contribution.
* **Slack:** Join Harbor's community for discussion and ask questions: [Cloud Native Computing Foundation](https://slack.cncf.io/), channel: [#harbor](https://cloud-native.slack.com/messages/harbor/), [#harbor-dev](https://cloud-native.slack.com/messages/harbor-dev/) and [#harbor-cli](https://cloud-native.slack.com/messages/harbor-cli/).

# License

This project is licensed under the Apache 2.0 License. See the [LICENSE](https://github.com/goharbor/harbor-cli/blob/main/LICENSE) file for details.

# Acknowledgements

This project is maintained by the Harbor community. We thank all our contributors and users for their support.

# ❤️ Show your support

For any questions or issues, please open an issue on our [GitHub Issues](https://github.com/goharbor/harbor-cli/issues) page.<br>
Give a ⭐ if this project helped you, Thank YOU!
