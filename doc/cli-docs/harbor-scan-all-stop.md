---
title: harbor scan all stop
weight: 40
---
## harbor scan-all stop

### Description

##### Stop scanning all artifacts

### Synopsis

Stop an ongoing vulnerability scan of all artifacts in Harbor.

This command halts the current scan-all operation that was either manually triggered 
or scheduled. When stopped, scans that are already in progress will complete, but no new artifacts will be scanned. The scan can be restarted later using the 'scan-all run' command.

Examples:
  # Stop the current scan-all operation
  harbor-cli scan-all stop

  # Stop and then check metrics to confirm
  harbor-cli scan-all stop && harbor-cli scan-all metrics

```sh
harbor scan-all stop [flags]
```

### Options

```sh
  -h, --help   help for stop
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor scan-all](harbor-scan-all.md)	 - Scan all artifacts

