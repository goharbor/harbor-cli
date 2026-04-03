---
title: harbor jobservice pools
weight: 45
---
## harbor jobservice pools

### Description

##### Manage worker pools (list available pools)

### Synopsis

List and manage worker pools for the Harbor job service.

Use 'list' to view all worker pools.

Examples:
  harbor jobservice pools list
  harbor jobservice pool list

### Options

```sh
  -h, --help   help for pools
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice](harbor-jobservice.md)	 - Manage Harbor job service (admin only)
* [harbor jobservice pools list](harbor-jobservice-pools-list.md)	 - List all worker pools

