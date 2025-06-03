---
title: harbor context delete
weight: 85
---
## harbor context delete

### Description

##### Delete (clear) a specific config item

### Synopsis

Clear the value of a specific CLI config item by setting it to its zero value.
Case-insensitive field lookup, but uses the canonical (Go) field name internally.
If you specify --name, that credential (rather than the "current" one) will be used.

```sh
harbor context delete <item> [flags]
```

### Examples

```sh

  # Clear the current credential's password
  harbor context delete credentials.password

  # Clear a specific credential's password using --name
  harbor context delete credentials.password --name admin@http://demo.goharbor.io

  # Clear the current credential
  harbor context delete --current

```

### Options

```sh
      --current       Remove current credentials from the config
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

* [harbor context](harbor-context.md)	 - Manage locally available contexts

