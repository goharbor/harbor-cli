---
title: harbor context update
weight: 10
---
## harbor context update

### Description

##### Set/update a specific config item

### Synopsis

Set/update the value of a specific CLI config item. 
Case-insensitive field lookup, but uses the canonical (Go) field name internally.
If you specify --name, that credential (rather than the "current" one) will be updated.

```sh
harbor context update <item> <value> [flags]
```

### Examples

```sh

  # Set/update the current credential's password
  harbor context update credentials.password myNewSecret

  # Set/update a credential's password by specifying the credential name
  harbor config update credentials.password myNewSecret --name admin@http://demo.goharbor.io

```

### Options

```sh
  -h, --help          help for update
  -n, --name string   Name of the credential to set fields on (default: the current credential)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor context](harbor-context.md)	 - Manage locally available contexts

