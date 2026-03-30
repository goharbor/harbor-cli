---
title: harbor jobservice jobs log
weight: 65
---
## harbor jobservice jobs log

### Description

##### View a job log (--job-id required)

### Synopsis

Display the log for a specific job by job ID.

```sh
harbor jobservice jobs log [flags]
```

### Examples

```sh
harbor jobservice jobs log --job-id abc123def456
```

### Options

```sh
  -h, --help            help for log
      --job-id string   Job ID to fetch log for (required)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice jobs](harbor-jobservice-jobs.md)	 - Manage job logs (view by job ID)

