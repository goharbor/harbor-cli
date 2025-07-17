---
title: harbor replication executions list
weight: 10
---
## harbor replication executions list

### Description

##### List replication executions

### Synopsis

List all replication executions for a given replication policy. If no policy ID is provided, it will prompt the user to select one interactively.

```sh
harbor replication executions list [flags]
```

### Examples

```sh
  harbor replication executions list 12345
  harbor replication executions list
```

### Options

```sh
  -h, --help            help for list
      --page int        Page number (default 1)
      --page-size int   Size of per page (0 to fetch all)
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

* [harbor replication executions](harbor-replication-executions.md)	 - Manage replication executions

