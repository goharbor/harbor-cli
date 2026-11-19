---
title: harbor jobservice jobs log
weight: 35
---
## harbor jobservice jobs log

### Description

##### View a job log (--job-id required)

### Synopsis

Display the log for a specific job by job ID.

Job logs contain detailed execution information including status updates, error messages, and processing details.
The --job-id flag is required to specify which job's log to retrieve.

Job IDs can be obtained from the 'harbor jobservice jobs list' command.

```sh
harbor jobservice jobs log [flags]
```

### Examples

```sh
View a specific job log:
  harbor jobservice jobs log --job-id abc123def456

View log with verbose output:
  harbor jobservice jobs log --job-id abc123def456 -v
```

### Options

```sh
  -h, --help            help for log
      --job-id string   Job ID to fetch log for (required)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice jobs](harbor-jobservice-jobs.md)	 - Manage job logs (view by job ID)

