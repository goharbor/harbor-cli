---
title: harbor artifact tags delete
weight: 70
---
## harbor artifact tags delete

### Description

##### Delete a tag of an artifact

```sh
harbor artifact tags delete [flags]
```

### Examples

```sh
harbor artifact tags delete <project>/<repository>/<reference> <tag>
```

### Options

```sh
  -h, --help   help for delete
```

### Options inherited from parent commands

```sh
      --config string          config file (default is $HOME/.harbor/config.yaml) (default "/home/user/.harbor/config.yaml")
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor artifact tags](harbor-artifact-tags.md)	 - Manage tags of an artifact

