---
title: harbor config set
weight: 85
---
## harbor config set

### Description

##### Set a specific config item

### Synopsis

Set the value of a specific CLI config item. 
Case-insensitive field lookup, but uses the canonical (Go) field name internally.
If you specify --name, that credential (rather than the "current" one) will be updated.

```sh
harbor config set <item> <value> [flags]
```

### Examples

```sh

  # Set the current credential's password
  harbor config set credentials.password myNewSecret

  # Set a credential's password by specifying the credential name
  harbor config set credentials.password myNewSecret --name harbor-cli@http://demo.goharbor.io

```

### Options

```sh
  -h, --help          help for set
  -n, --name string   Name of the credential to set fields on (default: the current credential)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor config](harbor-config.md)	 - Manage the config of the Harbor Cli

