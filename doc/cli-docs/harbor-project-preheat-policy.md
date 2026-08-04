---
title: harbor project preheat policy
weight: 90
---
## harbor project preheat policy

### Description

##### Manage preheat policies

### Synopsis

Manage P2P preheat policies under a project

### Examples

```sh
  harbor project preheat policy list [PROJECT_NAME]
  harbor project preheat policy list [PROJECT_ID] --id
  harbor project preheat policy create -f [CONFIG_FILE]
  harbor project preheat policy view [PROJECT_NAME] [POLICY_NAME]
  harbor project preheat policy start [PROJECT_NAME] [POLICY_NAME]
```

### Options

```sh
  -h, --help   help for policy
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project preheat](harbor-project-preheat.md)	 - Manage project preheat resources
* [harbor project preheat policy create](harbor-project-preheat-policy-create.md)	 - Create a preheat policy
* [harbor project preheat policy delete](harbor-project-preheat-policy-delete.md)	 - Delete a preheat policy
* [harbor project preheat policy list](harbor-project-preheat-policy-list.md)	 - List preheat policies under a project
* [harbor project preheat policy start](harbor-project-preheat-policy-start.md)	 - Manually trigger a preheat policy
* [harbor project preheat policy update](harbor-project-preheat-policy-update.md)	 - Update a preheat policy
* [harbor project preheat policy view](harbor-project-preheat-policy-view.md)	 - View details of a preheat policy

