---
title: harbor project preheat policy update
weight: 25
---
## harbor project preheat policy update

### Description

##### Update a preheat policy

### Synopsis

Update an existing P2P preheat policy under a project

```sh
harbor project preheat policy update [flags]
```

### Examples

```sh
  harbor project preheat policy update [PROJECT_NAME] [POLICY_NAME]
  harbor project preheat policy update [PROJECT_ID] [POLICY_NAME] --id
```

### Options

```sh
  -h, --help[=true|false]   help for update
      --id[=true|false]     Use project id instead of name
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -l, --log-format string      Output format for logging. One of: json|text (default "text")
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose[=true|false]   verbose output
```

### SEE ALSO

* [harbor project preheat policy](harbor-project-preheat-policy.md)	 - Manage preheat policies

