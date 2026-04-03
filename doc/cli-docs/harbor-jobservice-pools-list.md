---
title: harbor jobservice pools list
weight: 70
---
## harbor jobservice pools list

### Description

##### List all worker pools

### Synopsis

Display all worker pools with their details.

```sh
harbor jobservice pools list [flags]
```

### Examples

```sh
harbor jobservice pools list
harbor jobservice pool list
```

### Options

```sh
  -h, --help   help for list
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice pools](harbor-jobservice-pools.md)	 - Manage worker pools (list available pools)

