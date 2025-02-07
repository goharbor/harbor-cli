---
title: harbor artifact list
weight: 25
---
## harbor artifact list

### Description

##### list artifacts within a repository

```sh
harbor artifact list [flags]
```

### Options

```sh
  -h, --help            help for list
  -p, --page int        Page number (default 1)
  -n, --page-size int   Size of per page (default 10)
  -q, --query string    Query string to query resources
  -s, --sort string     Sort the resource list in ascending or descending order
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor artifact](harbor-artifact.md)	 - Manage artifacts

