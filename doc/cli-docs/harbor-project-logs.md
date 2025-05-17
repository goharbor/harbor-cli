---
title: harbor project logs
weight: 90
---
## harbor project logs

### Description

##### get project logs

```sh
harbor project logs [flags]
```

### Options

```sh
  -h, --help            help for logs
      --page int        Page number (default 1)
      --page-size int   Size of per page (default 10)
  -q, --query string    Query string to query resources
      --sort string     Sort the resource list in ascending or descending order
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project](harbor-project.md)	 - Manage projects and assign resources to them

