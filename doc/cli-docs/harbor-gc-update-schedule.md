---
title: harbor gc update schedule
weight: 50
---
## harbor gc update-schedule

### Description

##### Update automatic GC schedule

### Synopsis

Configure or update the automatic Garbage Collection schedule for the registry.

Available schedule types:
  - none:    Disable scheduled Garbage Collection
  - hourly:  Run GC every hour
  - daily:   Run GC once per day
  - weekly:  Run GC once per week
  - custom:  Define a custom schedule using a cron expression

For custom schedules, Harbor requires a 6-field cron expression in the format:
  seconds minutes hours day-of-month month day-of-week

Examples:
  # Disable automatic Garbage Collection
  harbor gc update-schedule none

  # Configure daily Garbage Collection deleting untagged artifacts
  harbor gc update-schedule daily --delete-untagged

  # Configure custom schedule (e.g. daily at 3:00 AM) in dry-run mode
  harbor gc update-schedule custom --cron "0 0 3 * * *" --dry-run

```sh
harbor gc update-schedule [schedule-type: none|hourly|daily|weekly|custom] [flags]
```

### Options

```sh
      --cron string       Cron expression for custom schedule (include in double quotes)
      --delete-untagged   Delete untagged artifacts
      --dry-run           Simulate the GC process without deleting actual blobs
  -h, --help              help for update-schedule
  -i, --interactive       Update GC schedule interactively
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor gc](harbor-gc.md)	 - Manage Garbage Collection in Harbor

