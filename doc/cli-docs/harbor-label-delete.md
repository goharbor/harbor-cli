---
title: harbor label delete
weight: 0
---
## harbor label delete

### Description

##### delete label

```sh
harbor label delete [flags]
```

### Examples

```sh
harbor label delete [labelname]
```

### Options

```sh
  -h, --help           help for delete
  -s, --scope string   default(global).'p' for project labels.Query scope of the label (default "g")
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor label](harbor-label.md)	 - Manage labels in Harbor

