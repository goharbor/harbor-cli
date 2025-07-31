---
title: harbor config update
weight: 35
---
## harbor config update

### Description

##### Update system configurations from local config file

### Synopsis

Update Harbor system configurations using the values stored in your local config file.
This will push the configurations from your local config file to the Harbor server.
Make sure to run 'harbor config get' first to populate the local config file with current configurations.

```sh
harbor config update [flags]
```

### Examples

```sh
harbor config update
```

### Options

```sh
  -h, --help   help for update
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor config](harbor-config.md)	 - Manage system configurations

