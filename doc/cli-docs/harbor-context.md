---
title: harbor context
weight: 30
---
## harbor context

### Description

##### Manage locally available contexts

### Synopsis

The context command allows you to manage configuration items of the Harbor CLI.
				You can add, get, or delete specific configuration items, as well as list all configuration items of the Harbor CLI.

### Examples

```sh
harbor context list
```

### Options

```sh
  -h, --help   help for context
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI
* [harbor context delete](harbor-context-delete.md)	 - Delete (clear) a specific config item
* [harbor context get](harbor-context-get.md)	 - Get a specific config item
* [harbor context list](harbor-context-list.md)	 - List contexts
* [harbor context switch](harbor-context-switch.md)	 - Switch to a new context
* [harbor context update](harbor-context-update.md)	 - Set/update a specific config item

