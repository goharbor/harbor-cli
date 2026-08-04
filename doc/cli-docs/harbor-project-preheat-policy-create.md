---
title: harbor project preheat policy create
weight: 50
---
## harbor project preheat policy create

### Description

##### Create a preheat policy

### Synopsis

Create a new P2P preheat policy under a project

```sh
harbor project preheat policy create [flags]
```

### Examples

```sh
  harbor project preheat policy create [PROJECT_NAME]
  harbor project preheat policy create [PROJECT_ID] --id
  harbor project preheat policy create -f [CONFIG_FILE]
```

### Options

```sh
  -h, --help                        help for create
      --id                          Use project id instead of name
  -f, --policy-config-file string   YAML/JSON file with preheat policy configuration
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project preheat policy](harbor-project-preheat-policy.md)	 - Manage preheat policies

