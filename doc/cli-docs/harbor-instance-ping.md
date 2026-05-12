---
title: harbor instance ping
weight: 25
---
## harbor instance ping

### Description

##### Ping preheat provider instance by name or id

### Synopsis

Ping a preheat provider instance to test its connectivity in Harbor. You can specify the instance
by name or ID directly as an argument. If no argument is provided, you will be prompted to select
an instance from a list of available instances.

```sh
harbor instance ping [NAME|ID] [flags]
```

### Examples

```sh
  harbor-cli instance ping my-instance
  harbor-cli instance ping 1 --id
```

### Options

```sh
  -h, --help   help for ping
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

