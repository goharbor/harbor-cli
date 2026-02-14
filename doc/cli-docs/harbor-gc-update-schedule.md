---
title: harbor gc update schedule
weight: 70
---
## harbor gc update-schedule

### Description

##### update-schedule [schedule-type: none|hourly|daily|weekly|custom]

### Synopsis

Configure or update the automatic GC schedule.

Available schedule types:
  - none:    Disable automatic GC
  - hourly:  Run GC every hour
  - daily:   Run GC once per day
  - weekly:  Run GC once per week
  - custom:  Define a custom schedule using a cron expression

```sh
harbor gc update-schedule [flags]
```

### Options

```sh
      --cron string   Cron expression for custom schedule (include the expression in double quotes)
  -h, --help          help for update-schedule
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor gc](harbor-gc.md)	 - Manage Garbage Collection

