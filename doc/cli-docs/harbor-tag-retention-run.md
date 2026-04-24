---
title: harbor tag retention run
weight: 80
---
## harbor tag retention run

### Description

##### run retention policy

### Synopsis

trigger retention execution for a project policy

```sh
harbor tag retention run [PROJECT_NAME] [flags]
```

### Options

```sh
      --dry-run            trigger dry-run execution without deleting artifacts
  -h, --help               help for run
      --retention-id int   retention policy ID
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor tag retention](harbor-tag-retention.md)	 - Manage retention policies

