---
title: harbor artifact labels remove
weight: 10
---
## harbor artifact labels remove

### Description

##### Remove a label of an artifact

```sh
harbor artifact labels remove [flags]
```

### Examples

```sh
harbor artifact labels remove <project>/<repository>/<reference> <labelName|labelID>
```

### Options

```sh
  -h, --help   help for remove
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor artifact labels](harbor-artifact-labels.md)	 - Manage labels of an artifact

