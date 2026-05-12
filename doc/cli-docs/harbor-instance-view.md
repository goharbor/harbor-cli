---
title: harbor instance view
weight: 40
---
## harbor instance view

### Description

##### get preheat provider instance by name or id

### Synopsis

Get detailed information about a preheat provider instance in Harbor. You can specify the instance
by name or ID directly as an argument. If no argument is provided, you will be prompted to select
an instance from a list of available instances.

```sh
harbor instance view [NAME|ID] [flags]
```

### Examples

```sh
  harbor-cli instance view my-instance
  harbor-cli instance view 1 --id
```

### Options

```sh
  -h, --help   help for view
      --id     Get instance by id
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor instance](harbor-instance.md)	 - Manage preheat provider instances in Harbor

