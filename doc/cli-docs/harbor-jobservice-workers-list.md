---
title: harbor jobservice workers list
weight: 10
---
## harbor jobservice workers list

### Description

##### List workers (supports --page and --page-size)

### Synopsis

List job service workers.

Pagination:
	- --page selects the 1-based page number
	- --page-size controls how many workers are shown per page

Examples:
  harbor jobservice workers list
  harbor jobservice workers list --pool default
	harbor jobservice workers list --page 2 --page-size 20

```sh
harbor jobservice workers list [flags]
```

### Options

```sh
  -h, --help            help for list
      --page int        Page number (default 1)
      --page-size int   Number of workers per page (default 20)
      --pool string     Worker pool ID to list workers from (default: all) (default "all")
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice workers](harbor-jobservice-workers.md)	 - Manage workers (list all/by pool, free, free-all)

