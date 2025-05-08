---
title: harbor config
weight: 30
---
## harbor config

### Description

##### Manage the config of the Harbor CLI

### Synopsis

The config command allows you to manage configurations of the Harbor CLI.
				You can add, get, or delete specific config item, as well as list all config items of the Harbor Cli

### Options

```sh
  -h, --help   help for config
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI
* [harbor config delete](harbor-config-delete.md)	 - Delete (clear) a specific config item
* [harbor config get](harbor-config-get.md)	 - Get a specific config item
* [harbor config list](harbor-config-list.md)	 - List config items
* [harbor config update](harbor-config-update.md)	 - Set/update a specific config item

