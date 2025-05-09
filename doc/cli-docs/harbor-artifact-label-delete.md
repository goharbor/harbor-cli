---
title: harbor artifact label delete
weight: 35
---
## harbor artifact label delete

### Description

##### del label to an artifact

### Synopsis

del label to artifact

```sh
harbor artifact label delete [flags]
```

### Examples

```sh
harbor artifact label del <project>/<repository>/<reference> <label name>
```

### Options

```sh
  -h, --help   help for delete
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor artifact label](harbor-artifact-label.md)	 - label command for artifacts

