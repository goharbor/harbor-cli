---
title: harbor jobservice workers
weight: 90
---
## harbor jobservice workers

### Description

##### Manage workers (list all/by pool, free, free-all)

### Synopsis

Manage job service workers using the job service API.

Use 'list' to view workers from all pools or a specific pool.
Use 'free' and 'free-all' to stop running jobs and release busy workers.

Examples:
  harbor jobservice workers list
  harbor jobservice workers list --pool all
  harbor jobservice workers list --pool default
	harbor jobservice workers list --page 2 --page-size 20
	harbor jobservice workers list default
	harbor jobservice worker list 72327cf790564e45b7c89a2d
  harbor jobservice workers free --job-id <job-id>
  harbor jobservice workers free-all

### Options

```sh
  -h, --help   help for workers
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice](harbor-jobservice.md)	 - Manage Harbor job service (admin only)
* [harbor jobservice workers free](harbor-jobservice-workers-free.md)	 - Free one worker (--job-id required)
* [harbor jobservice workers free-all](harbor-jobservice-workers-free-all.md)	 - Free all busy workers (job-id=all)
* [harbor jobservice workers list](harbor-jobservice-workers-list.md)	 - List workers (supports --page and --page-size)

