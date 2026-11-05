---
title: harbor project preheat execution stop
weight: 5
---
## harbor project preheat execution stop

### Description

##### Stop preheat execution

### Synopsis

Stop a specific P2P preheat execution of a policy under a project

```sh
harbor project preheat execution stop [PROJECT_NAME|ID] [POLICY_NAME] [EXECUTION_ID] [flags]
```

### Examples

```sh
  harbor-cli project preheat execution stop [PROJECT_NAME|ID] [POLICY_NAME] [EXECUTION_ID]
```

### Options

```sh
  -h, --help   help for stop
      --id     Use project ID instead of name
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -l, --log-format string      Output format for logging. One of: json|text (default "text")
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project preheat execution](harbor-project-preheat-execution.md)	 - Manage preheat executions

