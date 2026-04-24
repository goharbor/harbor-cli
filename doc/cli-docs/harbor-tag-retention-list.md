---
title: harbor tag retention list
weight: 65
---
## harbor tag retention list

### Description

##### display retention policy for a project

### Synopsis

retrieve and display retention policy configured for a project

```sh
harbor tag retention list [PROJECT_NAME] [flags]
```

### Options

```sh
  -h, --help               help for list
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

