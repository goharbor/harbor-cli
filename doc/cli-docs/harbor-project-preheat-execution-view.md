---
title: harbor project preheat execution view
weight: 50
---
## harbor project preheat execution view

### Description

##### View preheat execution details

### Synopsis

Get details of a specific P2P preheat execution under a project

```sh
harbor project preheat execution view [PROJECT_NAME|ID] [POLICY_NAME] [EXECUTION_ID] [flags]
```

### Examples

```sh
  harbor-cli project preheat execution view [NAME|ID] [POLICY_NAME] [EXECUTION_ID]
```

### Options

```sh
  -h, --help   help for view
      --id     Get preheat policy execution by project id
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project preheat execution](harbor-project-preheat-execution.md)	 - Manage preheat executions

