---
title: harbor artifact label add
weight: 25
---
## harbor artifact label add

### Description

##### add label to an artifact

### Synopsis

add label to artifact

```sh
harbor artifact label add [flags]
```

### Examples

```sh
harbor artifact label add <project>/<repository>/<reference> <label name>
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

* [harbor artifact label](harbor-artifact-label.md)	 - label command for artifacts

