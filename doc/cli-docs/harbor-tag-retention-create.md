---
title: harbor tag retention create
weight: 30
---
## harbor tag retention create

### Description

##### create retention policy

### Synopsis

create a retention policy for a project

```sh
harbor tag retention create [PROJECT_NAME] [flags]
```

### Examples

```sh

# Create a retention policy for a specific project
harbor tag retention create my-project

# Create a retention policy interactively
harbor tag retention create
```

### Options

```sh
  -h, --help   help for create
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor tag retention](harbor-tag-retention.md)	 - Manage retention policies

