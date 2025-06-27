---
title: harbor scan all run
weight: 80
---
## harbor scan-all run

### Description

##### Scan all artifacts now

### Synopsis

Initiate an immediate vulnerability scan of all artifacts in Harbor.

This command triggers a manual scan of all artifacts across all projects in your Harbor instance.
The scan will check for known vulnerabilities in container images using the configured scanner(s).

The scan runs asynchronously in the background. After initiating the scan, you can:
  - Check the scan progress with 'harbor-cli scan-all metrics'
  - View results through the Harbor UI 

Important considerations:
  - This operation can be resource intensive on large registries
  - Scanning many artifacts simultaneously may impact system performance
  - The time to complete depends on the number and size of artifacts
  - Only one scan-all operation can run at a time

Examples:
  # Start scanning all artifacts immediately
  harbor-cli scan-all run

  # Start scanning and monitor progress
  harbor-cli scan-all run && watch -n 0.2 harbor-cli scan-all metrics

The scan progress and results can be monitored through the metrics command
or through the Harbor web interface.

```sh
harbor scan-all run [flags]
```

### Options

```sh
  -h, --help   help for run
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor scan-all](harbor-scan-all.md)	 - Scan all artifacts

