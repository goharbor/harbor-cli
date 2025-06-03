---
title: harbor artifact label add
weight: 95
---
## harbor artifact label add

### Description

##### Attach a label to an artifact in a Harbor project repository

### Synopsis

Attach an existing label to a specific artifact identified by <project>/<repository>:<reference>.
You can specify the artifact and label directly as arguments, or interactively select them if arguments are omitted.

Examples:
  # Add a label to an artifact using project/repo:reference and label name
  harbor artifact label add myproject/myrepo@sha256:abcdef1234567890 dev

  # Prompt-based label selection for an artifact
  harbor artifact label add library/nginx:1.21

  # Fully interactive mode (prompt for everything)
  harbor artifact label add


```sh
harbor artifact label add [flags]
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

