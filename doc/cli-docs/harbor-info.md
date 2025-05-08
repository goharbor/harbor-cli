---
title: harbor info
weight: 65
---
## harbor info

### Description

##### Display detailed Harbor system, statistics, and CLI environment information

### Synopsis

The 'info' command retrieves and displays general information about the Harbor instance, 
including system metadata, storage statistics, and CLI environment details such as user identity, 
registry address, and CLI version.

The output can be formatted as table (default), JSON, or YAML using the '--output-format' flag.

```sh
harbor info [flags]
```

### Examples

```sh
  harbor info
  harbor info --output-format json
  harbor info -o yaml
```

### Options

```sh
  -h, --help   help for info
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI

