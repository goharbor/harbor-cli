---
title: harbor instance delete
weight: 35
---
## harbor instance delete

### Description

##### Delete a preheat provider instance by its name or ID

### Synopsis

Delete a preheat provider instance from Harbor. You can specify the instance name or ID directly as an argument.
If no argument is provided, you will be prompted to select an instance from a list of available instances.

```sh
harbor instance delete [flags]
```

### Examples

```sh
  harbor-cli instance delete my-instance
  harbor-cli instance delete 12345
```

### Options

```sh
  -h, --help     help for delete
  -i, --id int   ID of the instance to delete (default -1)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor instance](harbor-instance.md)	 - Manage preheat provider instances in Harbor

