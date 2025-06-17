---
title: Harbor CLI Config Management
weight: 25
---

# Harbor CLI Configuration Management

> **Note**  
> The Harbor CLI follows the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) for configuration and data storage by default.

## Introduction
Harbor CLI is a flexible command-line tool that lets you manage various Harbor environments with different credentials. Whether you need a production-ready setup or quick testing configurations, the CLI's hierarchical structure and XDG support help keep things organized.


## Understanding the Configuration Structure
The Harbor CLI can manage multiple credentials and keep track of which credential is currently active. This setup allows you to maintain separate contexts for different Harbor instances or user accounts without having to rewrite configuration files manually. While the Harbor CLI configuration file manages your credentials, passwords themselves are never stored in plain text. Instead, they are secured using the AES-GCM encryption described in the [Harbor CLI Encryption documentation](../cli-config).

### Example Configuration File
Below is a simplified example of a typical Harbor CLI configuration file:
```yaml
current-credential-name: example@demo-harbor
credentials:
  - name: example@demo-harbor
    username: example-user
    password: example-password
    serveraddress: https://demo.goharbor.io
```

In this configuration:
- **current-credential-name** references the active credential by name.  
- **credentials** holds one or more sets of user credentials, each following the same structure.

## Managing Multiple Credentials
If you need to work with multiple sets of credentials—such as development, staging, or production — Harbor CLI makes it easy to create and switch between them.

> **Note**: For more login command details please refer to the [login command reference](../cli-docs/harbor-login.md).
### Creating a New Credentials Entry
Use the `harbor login` command with the required arguments to store new credentials:
```bash
harbor login --name my-new-credential \
  --username myuser \
  --password mypass \
  https://my-harbor-instance.com
```
This adds a new entry to your credentials list, allowing you to manage different Harbor accounts from the same CLI.

### Switching Between Credentials
To switch to another credential set, run:
```bash
harbor login --name <name-of-credential>
```
The CLI will then set the specified credential as the active one, eliminating the need to manually edit your configuration files. This will overwrite the `current-credential-name`.


## Configuration Hierarchy (Highest to Lowest Priority)

1. **Explicit Config Flag**  
   Provide a custom config file at runtime using `--config`:
   ```bash
   harbor --config /path/to/custom/config.yaml artifact list
   ```

2. **Environment Variable**  
   Set a persistent configuration through the `HARBOR_CLI_CONFIG` environment variable:
   ```bash
   export HARBOR_CLI_CONFIG="$HOME/.custom/harbor-config.yaml"
   harbor artifact list  # Uses the environment-specified config
   ```

3. **XDG Default Paths**  
   Automatically discover configuration in the following order:
   ```bash
   ${XDG_CONFIG_HOME}/harbor-cli/config.yaml  # If XDG_CONFIG_HOME is set
   ~/.config/harbor-cli/config.yaml           # Fallback default
   ```

## Data Storage Management
### Data File Location

- **Primary Path**: `$XDG_DATA_HOME/harbor-cli/data.yaml`
- **Fallback Path**: `$HOME/.local/share/harbor-cli/data.yaml`

> **Important**  
> The data file automatically tracks the last-used configuration file path

## Configuration Precedence Summary
| Priority | Method                     | Example                               |
|----------|----------------------------|---------------------------------------|
| 1        | --config flag             | harbor --config ./test.yaml ...       |
| 2        | HARBOR_CLI_CONFIG env var | export HARBOR_CLI_CONFIG=...          |
| 3        | XDG Default Locations     | ~/.config/harbor-cli/config.yaml      |

## Practical Usage Examples
### Scenario 1: Temporary Config Override
```bash
harbor --config ./experimental.yaml project create "new-project"
```

### Scenario 2: Persistent Environment-based Config
```bash
echo 'export HARBOR_CLI_CONFIG="$HOME/work/configs/prod-harbor.yaml"' >> ~/.zshrc
source ~/.zshrc
harbor config list  # Uses production config
```

### Scenario 3: Reset to Default Configuration
```bash
unset HARBOR_CLI_CONFIG
harbor config delete --current  # Deletes current context
```
