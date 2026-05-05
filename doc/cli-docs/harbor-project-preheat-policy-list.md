---
title: harbor project preheat policy list
weight: 20
---
## harbor project preheat policy list

### Description

##### List preheat policies under a project by name or id

### Synopsis

List project-scoped P2P preheat policies in Harbor

```sh
harbor project preheat policy list [NAME|ID] [flags]
```

### Options

```sh
  -h, --help            help for list
      --id              Get preheat policies by project id
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

* [harbor project preheat policy](harbor-project-preheat-policy.md)	 - Manage preheat policies

