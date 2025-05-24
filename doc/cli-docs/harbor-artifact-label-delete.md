---
title: harbor artifact label delete
weight: 90
---
## harbor artifact label delete

### Description

##### Detach a label from an artifact in a Harbor project repository

### Synopsis

Remove an existing label from a specific artifact identified by <project>/<repository>:<reference>.
You can provide the artifact and label name as arguments, or choose them interactively if not specified.

Examples:
  # Remove a label by specifying artifact and label name
  harbor artifact label delete library/nginx:latest stable

  # Prompt-based label selection for a specific artifact
  harbor artifact label del library/nginx:1.21

  # Fully interactive mode (prompt for project, repo, reference, and label)
  harbor artifact label delete

  # Remove a label from an artifact identified by digest
  harbor artifact label del myproject/myrepo@sha256:abcdef1234567890 qa-label

```sh
harbor artifact label delete [flags]
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

