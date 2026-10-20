---
title: harbor jobservice
weight: 75
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
  -l, --log-format string      Output format for logging. One of: json|text (default "text")
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI
* [harbor jobservice queues](harbor-jobservice-queues.md)	 - Manage job queues (list, stop, pause, resume)

