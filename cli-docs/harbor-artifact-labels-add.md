---
title: harbor artifact labels add
weight: 95
---
## harbor artifact labels add

### Description

##### Add a label of an artifact

```sh
harbor artifact labels add [flags]
```

### Examples

```sh
harbor artifact labels add <project>/<repository>/<reference> <labelName|labelID>
```

### Options

```sh
  -h, --help   help for add
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor artifact labels](harbor-artifact-labels.md)	 - Manage labels of an artifact

