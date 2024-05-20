
![harbor-3](https://github.com/goharbor/harbor-cli/assets/70086051/835ab686-1cce-4ac7-bc57-05a35c2b73cc)

**Welcome to the Harbor CLI project! This powerful command-line tool facilitates seamless interaction with the Harbor container registry. It simplifies various tasks such as creating, updating, and managing projects, registries, and other resources in Harbor.**

# Project Scope 🧪

The Harbor CLI is designed to enhance your interaction with the Harbor container registry. Built on Golang, it offers a user-friendly interface to perform various tasks related to projects, registries, and more. Whether you're creating, updating, or managing resources, the Harbor CLI streamlines your workflow efficiently.

# **Project Features** 🤯

 🔹 Get details about projects, registries, repositories and more <br>
 🔹 Create new projects, registries, and other resources <br>
 🔹 Delete projects, registries, and other resources <br>
 🔹 Run commands with various flags for enhanced functionality <br>
 🔹 More features coming soon... 🚧

# Demo Screenshot

![demo-1](https://github.com/goharbor/harbor-cli/assets/70086051/63b5f6b7-319b-4c05-968a-59489f7fdd35)

![demo-2](https://github.com/goharbor/harbor-cli/assets/70086051/00afaa16-41c4-460d-8ec1-7b06b02bd80c)

# Supported Platforms

Platform | Status
--|--
Linux | ✅
macOS | ✅
Windows | ✅

# Installation

## Build From Source
```bash
git clone https://github.com/goharbor/harbor-cli.git
cd harbor-cli/cmd/harbor
go build .
sudo mv harbor /usr/local/bin/
```
## Linux and MacOS
 use `amd64/arm64` as per your system architecture
```bash
## Linux
tar -xzf harbor_0.0.1_linux_amd64.tar.gz
cd harbor_0.0.1_linux_amd64
sudo mv harbor /usr/local/bin/

## MacOS
tar -xzf harbor_0.0.1_darwin_amd64.tar.gz
cd harbor_0.0.1_darwin_amd64
sudo mv harbor /usr/local/bin/
```
## Windows
 - Download `harbor_0.0.1_windows_amd64.zip` and Extract it.
 - To easily use the harbor-cli from the command line, add the directory containing the `harbor.exe` to your system PATH.
 - In the Edit Environment Variable window, click on "New" and add the path to the directory where `harbor.exe` is located (e.g., `C:\path\to\harbor_0.0.1_windows_amd64`).

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
