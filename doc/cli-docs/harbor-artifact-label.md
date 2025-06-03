---
title: harbor artifact label
weight: 70
---
## harbor artifact label

### Description

##### label command for artifacts

### Synopsis

label command for artifact

### Examples

```sh
harbor artifact label add <project>/<repository>/<reference> <label name>
harbor artifact label del <project>/<repository>/<reference> <label name>
		
```

### Options

```sh
  -h, --help   help for label
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor artifact](harbor-artifact.md)	 - Manage artifacts
* [harbor artifact label add](harbor-artifact-label-add.md)	 - Attach a label to an artifact in a Harbor project repository
* [harbor artifact label delete](harbor-artifact-label-delete.md)	 - Detach a label from an artifact in a Harbor project repository
* [harbor artifact label list](harbor-artifact-label-list.md)	 - Display labels attached to a specific artifact

