---
title: harbor gc
weight: 75
---
## harbor gc

### Description

##### Manage Garbage Collection

### Synopsis

Manage Garbage Collection in Harbor (schedule, history, logs)

### Options

```sh
  -h, --help   help for gc
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI
* [harbor gc list](harbor-gc-list.md)	 - List GC history
* [harbor gc log](harbor-gc-log.md)	 - Get GC job log
* [harbor gc run](harbor-gc-run.md)	 - Run Garbage Collection manually
* [harbor gc schedule](harbor-gc-schedule.md)	 - Display the GC schedule
* [harbor gc update-schedule](harbor-gc-update-schedule.md)	 - update-schedule [schedule-type: none|hourly|daily|weekly|custom]

