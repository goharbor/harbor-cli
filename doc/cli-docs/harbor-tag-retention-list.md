---
title: harbor tag retention list
weight: 0
---
## harbor tag retention list

### Description

##### List tag retention rules of a project

### Synopsis

Retrieve and display the tag retention rules for a specific project in Harbor.

Tag retention rules define policies for automatically keeping or deleting image tags 
within a project. Using this command, you can view the currently configured 
retention rules.

Usage:
  - Specify the project **either by name or by ID**, but not both.
  - If neither is provided, you will be prompted to select a project.
  - The rules will be displayed in a formatted output.

Examples:
  # List retention rules using project name
  harbor tag retention list --project-name my-project

  # List retention rules using project ID
  harbor tag retention list --project-id 42

  # List retention rules in JSON format
  harbor tag retention list --project-name my-project --output-format json

```sh
harbor tag retention list [flags]
```

### Options

```sh
  -h, --help                  help for list
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

* [harbor tag retention](harbor-tag-retention.md)	 - Manage tag retention policies in the project

