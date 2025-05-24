---
title: harbor artifact label list
weight: 30
---
## harbor artifact label list

### Description

##### Display labels attached to a specific artifact

### Synopsis

This command lists all labels currently associated with a specific artifact in a Harbor project repository.
You can provide the artifact reference in the format <project>/<repository>:<reference> (where reference is either a tag or a digest).
If the reference is not provided as an argument, the command will prompt you to select the project, repository, and artifact.

Supports output formatting such as JSON or YAML using the --output (-o) flag.

```sh
harbor artifact label list [flags]
```

### Examples

```sh
  # List labels for a tagged artifact
  harbor artifact label list library/nginx:latest

  # List labels for an artifact by digest
  harbor artifact label list myproject/myrepo@sha256:abc123...

  # Prompt-based interactive selection of artifact
  harbor artifact label list

  # Output in JSON format
  harbor artifact label list library/nginx:1.21 -o json
```

### Options

```sh
  -h, --help   help for list
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor artifact label](harbor-artifact-label.md)	 - label command for artifacts

