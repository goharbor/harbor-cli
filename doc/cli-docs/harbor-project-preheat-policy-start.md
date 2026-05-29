---
title: harbor project preheat policy start
weight: 55
---
## harbor project preheat policy start

### Description

##### Manually trigger a preheat policy

### Synopsis

Manually trigger a specific P2P preheat policy under a project

```sh
harbor project preheat policy start [flags]
```

### Examples

```sh
  harbor project preheat policy start [PROJECT_NAME] [POLICY_NAME]
  harbor project preheat policy start [PROJECT_ID] [POLICY_NAME] --id
```

### Options

```sh
  -h, --help   help for start
      --id     Manually trigger preheat policy by project id
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project preheat policy](harbor-project-preheat-policy.md)	 - Manage preheat policies

