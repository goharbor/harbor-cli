---
title: harbor jobservice
weight: 80
---
## harbor jobservice

### Description

##### Manage Harbor job service (admin only)

### Synopsis

Manage Harbor job service components including worker pools, job queues, schedules, and job logs.
This requires system admin privileges.

Use "harbor jobservice [command] --help" for detailed examples and flags per subcommand.

### Options

```sh
  -h, --help   help for jobservice
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI
* [harbor jobservice jobs](harbor-jobservice-jobs.md)	 - Manage job logs (view by job ID)
* [harbor jobservice pools](harbor-jobservice-pools.md)	 - Manage worker pools (list available pools)
* [harbor jobservice queues](harbor-jobservice-queues.md)	 - Manage job queues (list, stop, pause, resume)
* [harbor jobservice schedules](harbor-jobservice-schedules.md)	 - Manage schedules (list/status/pause-all/resume-all)
* [harbor jobservice workers](harbor-jobservice-workers.md)	 - Manage workers (list all/by pool, free, free-all)

