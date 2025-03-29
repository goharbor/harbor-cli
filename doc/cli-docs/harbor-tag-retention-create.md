---
title: harbor tag retention create
weight: 50
---
## harbor tag retention create

### Description

##### Create a tag retention rule in a project

### Synopsis

Create a tag retention rule for a project in Harbor to manage the lifecycle of image tags.

Tag retention rules help users automatically retain or delete specific tags based on 
defined criteria, reducing storage usage and improving repository maintenance.

⚠️ A user can create **up to 15 tag retention rules per project**.

```sh
harbor tag retention create [flags]
```

### Examples

```sh
  # Retain tags matching 'release-*' at the project level
  harbor tag retention create --level project --action retain --taglist release-*

  # Delete untagged images at the repository level
  harbor retention create --level repository --action delete --tagdecoration untagged
```

### Options

```sh
      --action string           Action to perform: 'retain' or 'delete' (default "retain")
      --algorithm string        Rule combination method: 'or' or 'and' (default "or")
  -h, --help                    help for create
      --level string            Scope of the retention policy: 'project' or 'repository' (default "project")
      --repodecoration string   Apply or exclude repositories from the rule
      --repolist string         Comma-separated list of repositories to apply/exclude
      --tagdecoration string    Apply or exclude specific tags from the rule
      --taglist string          Comma-separated list of tags to apply/exclude
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor tag retention](harbor-tag-retention.md)	 - Manage tag retention policies in the project

