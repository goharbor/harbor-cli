---
title: harbor config delete
weight: 85
---
## harbor config delete

### Description

##### Delete (clear) a specific config item

### Synopsis

Clear the value of a specific CLI config item by setting it to its zero value.
Case-insensitive field lookup, but uses the canonical (Go) field name internally.
If you specify --name, that credential (rather than the "current" one) will be used.

```sh
harbor config delete <item> [flags]
```

### Examples

```sh

  # Clear the current credential's password
  harbor config delete credentials.password

  # Clear a specific credential's password using --name
  harbor config delete credentials.password --name harbor-cli@http://demo.goharbor.io

```

### Options

```sh
  -h, --help          help for delete
  -n, --name string   Name of the credential to delete fields from (default: the current credential)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor config](harbor-config.md)	 - Manage the config of the Harbor Cli

