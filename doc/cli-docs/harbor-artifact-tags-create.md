---
title: harbor artifact tags create
weight: 65
---
## harbor artifact tags create

### Description

##### Create a tag of an artifact

```sh
harbor artifact tags create [flags]
```

### Examples

```sh
harbor artifact tags create <project>/<repository>/<reference> <tag>
```

### Options

```sh
  -h, --help   help for create
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor artifact tags](harbor-artifact-tags.md)	 - Manage tags of an artifact

