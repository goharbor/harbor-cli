---
title: harbor context get
weight: 85
---
## harbor context get

### Description

##### Get a specific config item

### Synopsis

Get the value of a specific CLI config item.
If you specify --name, that credential (rather than the "current" one) will be used.

```sh
harbor context get <item> [flags]
```

### Examples

```sh

  # Get the current credential's username
  harbor context get credentials.username

  # Get a credential's username by specifying the credential name
  harbor config get credentials.username --name admin@http://demo.goharbor.io

```

### Options

```sh
  -h, --help          help for get
  -n, --name string   Name of the credential to get fields from (default: the current credential)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor context](harbor-context.md)	 - Manage locally available contexts

