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
  -h, --help            help for list
      --page int        Page number (default 1)
      --page-size int   Size of per page (default 20)
  -i, --projectid int   project ID when query project labels (default 1)
  -q, --query string    Query string to query resources
  -s, --scope string    default(global).'p' for project labels.Query scope of the label (default "g")
      --sort string     Sort the label list in ascending or descending order
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor label](harbor-label.md)	 - Manage labels in Harbor

