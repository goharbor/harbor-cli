---
title: harbor jobservice queues
weight: 35
---
## harbor jobservice queues

### Description

##### Manage job queues (list, stop, pause, resume)

### Synopsis

List job queues and perform actions on them (stop/pause/resume).

### Options

```sh
  -h, --help   help for queues
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice](harbor-jobservice.md)	 - Manage Harbor job service (admin only)
* [harbor jobservice queues list](harbor-jobservice-queues-list.md)	 - List all job queues
* [harbor jobservice queues pause](harbor-jobservice-queues-pause.md)	 - Pause queue(s) (--type or --interactive)
* [harbor jobservice queues resume](harbor-jobservice-queues-resume.md)	 - Resume queue(s) (--type or --interactive)
* [harbor jobservice queues stop](harbor-jobservice-queues-stop.md)	 - Stop queue(s) (--type or --interactive)

