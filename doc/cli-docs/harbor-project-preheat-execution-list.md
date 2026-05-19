---
title: harbor project preheat execution list
weight: 70
---
## harbor project preheat execution list

### Description

##### List preheat executions

### Synopsis

List preheat executions under a project

```sh
harbor project preheat execution list [flags]
```

### Examples

```sh
  harbor-cli project preheat execution list [NAME|ID] [POLICY_NAME]
```

### Options

```sh
  -h, --help            help for list
      --id              Get preheat executions by project id
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

* [harbor project preheat execution](harbor-project-preheat-execution.md)	 - Manage preheat executions

