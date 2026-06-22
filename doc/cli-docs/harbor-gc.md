---
title: harbor gc
weight: 85
---
## harbor gc

### Description

##### Manage Garbage Collection in Harbor

### Synopsis

Use this command to manage registry-wide Garbage Collection (GC) in your Harbor instance.

Garbage Collection cleans up deleted or orphaned blobs/tags in the registry to free up storage space.
This command supports listing execution history, viewing logs, showing schedule configuration, stopping running jobs, and triggering manual runs.

### Examples

```sh
  # View Garbage Collection execution history
  harbor gc history

  # Get the current Garbage Collection schedule
  harbor gc schedule

  # Trigger Garbage Collection run immediately
  harbor gc trigger --delete-untagged --dry-run=false

  # View execution logs for a GC run
  harbor gc log 12

  # Stop a running Garbage Collection run
  harbor gc stop 12

  # Update the automatic Garbage Collection schedule
  harbor gc update-schedule daily --delete-untagged
```

### Options

```sh
  -h, --help   help for gc
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI
* [harbor gc history](harbor-gc-history.md)	 - Get GC execution history
* [harbor gc log](harbor-gc-log.md)	 - Get GC execution log
* [harbor gc schedule](harbor-gc-schedule.md)	 - Get the current GC schedule
* [harbor gc stop](harbor-gc-stop.md)	 - Stop a running GC execution
* [harbor gc trigger](harbor-gc-trigger.md)	 - Trigger Garbage Collection immediately
* [harbor gc update-schedule](harbor-gc-update-schedule.md)	 - Update automatic GC schedule

