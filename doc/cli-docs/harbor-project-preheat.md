---
title: harbor project preheat
weight: 15
---
## harbor project preheat

### Description

##### Manage project preheat resources

### Synopsis

Manage project-scoped P2P preheat policies, executions, and tasks in Harbor

### Examples

```sh
  harbor project preheat policy list [PROJECT_NAME]
  harbor project preheat policy list [PROJECT_ID] --id
  harbor project preheat policy create -f [CONFIG_FILE]
  harbor project preheat policy start [PROJECT_NAME] [POLICY_NAME]
```

### Options

```sh
  -h, --help   help for preheat
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project](harbor-project.md)	 - Manage projects and assign resources to them
* [harbor project preheat execution](harbor-project-preheat-execution.md)	 - Manage preheat executions
* [harbor project preheat policy](harbor-project-preheat-policy.md)	 - Manage preheat policies
* [harbor project preheat task](harbor-project-preheat-task.md)	 - Manage preheat tasks

