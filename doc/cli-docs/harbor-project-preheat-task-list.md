---
title: harbor project preheat task list
weight: 65
---
## harbor project preheat task list

### Description

##### List preheat tasks

### Synopsis

List all tasks for a specific P2P preheat execution under a project

```sh
harbor project preheat task list [PROJECT_NAME|ID] [POLICY_NAME] [EXECUTION_ID] [flags]
```

### Examples

```sh
  harbor-cli project preheat task list [PROJECT_NAME|ID] [POLICY_NAME] [EXECUTION_ID]
```

### Options

```sh
  -h, --help            help for list
      --id              Get preheat tasks by project id
      --page int        Page number (default 1)
      --page-size int   Size of per page (default 10)
  -q, --query string    Query string to query resources
      --sort string     Sort the resource list in ascending or descending order
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project preheat task](harbor-project-preheat-task.md)	 - Manage preheat tasks

