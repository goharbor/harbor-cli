---
title: harbor gc log
weight: 70
---
## harbor gc log

### Description

##### Get GC job log

### Synopsis

Get the log of a specific GC (Garbage Collection) job.

If no GC job ID is provided via the --id flag, an interactive selector
will be displayed to choose from available GC jobs.

Examples:
  # Get GC log by specifying the job ID
  harbor gc log --id 42

  # Get GC log interactively (select from list)
  harbor gc log

```sh
harbor gc log [flags]
```

### Options

```sh
  -h, --help     help for log
      --id int   ID of the GC job to get logs for
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor gc](harbor-gc.md)	 - Manage Garbage Collection

