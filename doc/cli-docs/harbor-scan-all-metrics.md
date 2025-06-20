---
title: harbor scan all metrics
weight: 50
---
## harbor scan-all metrics

### Description

##### Get the metrics of the latest scan all process

### Synopsis

Display comprehensive metrics about the most recent vulnerability scan execution.

This command retrieves and displays detailed statistics about the most recent scan all 
process in Harbor, including:

- Running: Number of currently running scan tasks
- Success: Number of successfully completed scan tasks
- Error: Number of failed scan tasks
- Completed: Total number of completed scan tasks
- Total: Total number of scan tasks
- Ongoing: Whether the scan is still in progress
- Trigger: What triggered the scan (Manual, Scheduled, etc.)

The metrics provide visibility into the progress and results of vulnerability scanning across your Harbor registry.

Examples:
  # Get metrics for the latest scan
  harbor-cli scan-all metrics

  # Get metrics for the latest scheduled scan
  harbor-cli scan-all metrics --scheduled

  # Display metrics in JSON format
  harbor-cli scan-all metrics --output-format json

```sh
harbor scan-all metrics [flags]
```

### Options

```sh
  -h, --help        help for metrics
  -s, --scheduled   Get the metrics of the latest scheduled scan all process
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor scan-all](harbor-scan-all.md)	 - Scan all artifacts

