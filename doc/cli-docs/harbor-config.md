---
title: harbor config 
weight: 15
---
## harbor config

### Description

##### Manage the config of the Harbor cli

### Synopsis

#### Config
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

#### Data 
  - Data paths are determined by the `XDG_DATA_HOME` environment variable.
  - If `XDG_DATA_HOME` is not set, it defaults to `$HOME/.local/share/harbor-cli/data.yaml`.
  - The data file always contains the path of the latest config used.

### Options

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -h, --help                   help for harbor
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor artifact](harbor-artifact.md)	 - Manage artifacts
* [harbor login](harbor-login.md)	 - Log in to Harbor registry
* [harbor project](harbor-project.md)	 - Manage projects and assign resources to them
* [harbor registry](harbor-registry.md)	 - Manage registries
* [harbor repo](harbor-repo.md)	 - Manage repositories
* [harbor user](harbor-user.md)	 - Manage users
* [harbor version](harbor-version.md)	 - Version of Harbor CLI
