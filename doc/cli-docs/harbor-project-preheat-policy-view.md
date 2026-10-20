---
title: harbor project preheat policy view
weight: 30
---
## harbor project preheat policy view

### Description

##### View details of a preheat policy

### Synopsis

Get details of a specific P2P preheat policy under a project

```sh
harbor project preheat policy view [flags]
```

### Examples

```sh
  harbor project preheat policy view [PROJECT_NAME] [POLICY_NAME]
  harbor project preheat policy view [PROJECT_ID] [POLICY_NAME] --id
```

### Options

```sh
  -h, --help   help for view
      --id     Get preheat policy by project id
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -l, --log-format string      Output format for logging. One of: json|text (default "text")
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project preheat policy](harbor-project-preheat-policy.md)	 - Manage preheat policies

