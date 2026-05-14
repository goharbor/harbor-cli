---
title: harbor project preheat task log
weight: 60
---
## harbor project preheat task log

### Description

##### Get preheat task log

### Synopsis

Get the log for a specific P2P preheat task under a project

```sh
harbor project preheat task log [PROJECT_NAME|ID] [POLICY_NAME] [EXECUTION_ID] [TASK_ID] [flags]
```

### Examples

```sh
  harbor-cli project preheat task log [PROJECT_NAME|ID] [POLICY_NAME] [EXECUTION_ID] [TASK_ID]
```

### Options

```sh
  -h, --help   help for log
      --id     Use project ID instead of name
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project preheat task](harbor-project-preheat-task.md)	 - Manage preheat tasks

