---
title: harbor jobservice workers list
weight: 95
---
## harbor jobservice workers list

### Description

##### List workers (all pools by default; use --pool for one pool)

### Synopsis

List job service workers.

Supported listing modes:
	- All workers (default): no POOL_ID or --pool all
	- Specific pool workers: provide [POOL_ID] or --pool <pool-id>
	- Compatibility mode: --pool-all (same as --pool all)

Examples:
  harbor jobservice workers list
  harbor jobservice workers list --pool all
  harbor jobservice workers list --pool default
	harbor jobservice workers list default
	harbor jobservice worker list 72327cf790564e45b7c89a2d

```sh
harbor jobservice workers list [POOL_ID] [flags]
```

### Options

```sh
      --all           List workers from all pools
  -h, --help          help for list
      --pool string   Worker pool ID (use 'all' for all pools)
      --pool-all      List workers from all pools (compatibility alias for --pool all)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice workers](harbor-jobservice-workers.md)	 - Manage workers (list all/by pool, free, free-all)

