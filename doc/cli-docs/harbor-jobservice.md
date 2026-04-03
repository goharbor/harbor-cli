---
title: harbor jobservice
weight: 15
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
* [harbor jobservice schedules](harbor-jobservice-schedules.md)	 - Manage schedules (list/status/pause-all/resume-all)

