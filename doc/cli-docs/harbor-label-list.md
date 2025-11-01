---
title: harbor label list
weight: 15
---
## harbor label list

### Description

##### list labels

```sh
harbor label list [flags]
```

### Options

```sh
      --fuzzy strings    Fuzzy match filter (key=value)
      --global           whether to list global or project scope labels. (default scope is global)
  -h, --help             help for list
      --match strings    exact match filter (key=value)
      --page int         Page number (default 1)
      --page-size int    Size of per page (default 20)
  -p, --project string   project name when query project labels
  -i, --project-id int   project ID when query project labels
  -q, --query string     Query string to query resources
      --range strings    range filter (key=min~max)
      --sort string      Sort the label list in ascending or descending order
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor label](harbor-label.md)	 - Manage labels in Harbor

