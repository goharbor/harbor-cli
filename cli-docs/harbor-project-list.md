---
title: harbor project list
weight: 40
---
## harbor project list

### Description

##### list project

```sh
harbor project list [flags]
```

### Options

```sh
  -h, --help            help for list
      --name string     Name of the project
      --page int        Page number (default 1)
      --page-size int   Size of per page (default 10)
      --private         Show only private projects
      --public          Show only public projects
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

