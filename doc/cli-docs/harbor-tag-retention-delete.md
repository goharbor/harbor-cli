---
title: harbor tag retention delete
weight: 15
---
## harbor tag retention delete

### Description

##### Delete a tag retention rule for a project

### Synopsis

Delete an existing tag retention rule from a project.

Usage:
  - You can specify the project either by name or by ID, but not both.
  - If neither is provided, you will be prompted to select a project.
  - The command retrieves the retention policy ID and deletes it.

Examples:
  # Delete retention rule using project name
  harbor tag retention delete --project-name my-project

  # Delete retention rule using project ID
  harbor tag retention delete --project-id 42

```sh
harbor tag retention delete [flags]
```

### Options

```sh
  -h, --help                  help for delete
  -i, --project-id int        Project ID (default -1)
  -p, --project-name string   Project name
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor tag retention](harbor-tag-retention.md)	 - Manage tag retention rules in the project

