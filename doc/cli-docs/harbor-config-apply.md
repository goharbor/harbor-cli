---
title: harbor config apply
weight: 85
---
## harbor config apply

### Description

##### Update system configurations from local config file

### Synopsis

Update Harbor system configurations using the values stored in your local config file.
This will push the configurations from your local config file to the Harbor server.
Make sure to run 'harbor config get' first to populate the local config file with current configurations. Alternatively, you can specify a custom configuration file using the --configurations-file flag. This does not have to be a complete configuration file, only the fields you want to update need to be present under the 'configurations' key. Credentials for the Harbor server can be configured in the local config file or through environment variables or global config flags.

```sh
harbor config apply [flags]
```

### Examples

```sh
harbor config apply
```

### Options

```sh
  -f, --configurations-file string   Harbor configurations file to apply.
  -h, --help                         help for apply
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor config](harbor-config.md)	 - Manage system configurations

