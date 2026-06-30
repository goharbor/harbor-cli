---
title: harbor project preheat task
weight: 5
---
## harbor project preheat task

### Description

##### Manage preheat tasks

### Synopsis

Manage related tasks for the given preheat execution

### Examples

```sh
  harbor-cli project preheat task list [PROJECT_NAME|ID] [POLICY_NAME] [EXECUTION_ID]
```

### Options

```sh
  -h, --help   help for task
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project preheat](harbor-project-preheat.md)	 - Manage project preheat resources
* [harbor project preheat task list](harbor-project-preheat-task-list.md)	 - List preheat tasks
* [harbor project preheat task log](harbor-project-preheat-task-log.md)	 - Get preheat task log

