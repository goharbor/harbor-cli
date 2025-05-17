---
title: harbor repo list
weight: 150
---
## harbor repo list

### Description

##### list repositories within a project

### Synopsis

Get information of all repositories in a project

```sh
harbor repo list [flags]
```

### Examples

```sh
  harbor repo list <project_name>
```

### Options

```sh
  -h, --help            help for list
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

* [harbor repo](harbor-repo.md)	 - Manage repositories

