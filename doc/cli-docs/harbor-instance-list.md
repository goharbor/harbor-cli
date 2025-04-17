---
title: harbor instance list
weight: 95
---
## harbor instance list

### Description

##### List all preheat provider instances in Harbor

### Synopsis

List all preheat provider instances registered in Harbor. You can paginate the results, 
filter them using a query string, and sort them in ascending or descending order. 
This command provides an easy way to view all instances along with their details.

```sh
harbor instance list [flags]
```

### Examples

```sh
  harbor-cli instance list --page 1 --page-size 10
  harbor-cli instance list --query "name=my-instance" --sort "asc"
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

* [harbor instance](harbor-instance.md)	 - Manage preheat provider instances in Harbor

