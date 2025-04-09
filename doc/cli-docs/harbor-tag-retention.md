---
title: harbor tag retention
weight: 95
---
## harbor tag retention

### Description

##### Manage tag retention rules in the project

### Synopsis

Manage tag retention rules in the project in Harbor.
		
The 'retention' command allows users to create, list, and delete tag retention rules 
within a project. Tag retention rules help in managing and controlling the lifecycle 
of tags by defining rules for automatic cleanup and retention.

A user can only create up to 15 tag retention rules per project.

### Examples

```sh
  harbor tag retention create    # Create a new tag retention rule
  harbor tag retention list      # List all tag retention rules in the project
  harbor tag retention delete    # Delete a specific tag retention rules
```

### Options

```sh
  -h, --help   help for retention
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor tag](harbor-tag.md)	 - Manage tag rules in Harbor registry
* [harbor tag retention create](harbor-tag-retention-create.md)	 - Create a tag retention rule in a project
* [harbor tag retention delete](harbor-tag-retention-delete.md)	 - Delete a tag retention rule for a project
* [harbor tag retention list](harbor-tag-retention-list.md)	 - List tag retention rules of a project

