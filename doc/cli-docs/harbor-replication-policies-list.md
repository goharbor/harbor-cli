---
title: harbor replication policies list
weight: 85
---
## harbor replication policies list

### Description

##### List replication policies

```sh
harbor replication policies list [flags]
```

### Options

```sh
      --fuzzy strings   Fuzzy match filter (key=value)
  -h, --help            help for list
      --match strings   exact match filter (key=value)
      --page int        Page number (default 1)
      --page-size int   Size of per page (0 to fetch all)
  -q, --query string    Query string to query resources
      --range strings   range filter (key=min~max)
      --sort string     Sort the resource list in ascending or descending order
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor replication policies](harbor-replication-policies.md)	 - Manage replication policies

