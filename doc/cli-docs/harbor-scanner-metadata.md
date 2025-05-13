---
title: harbor scanner metadata
weight: 60
---
## harbor scanner metadata

### Description

##### Retrieve metadata for a specific scanner

### Synopsis

Retrieve detailed metadata for a specified scanner integration in Harbor.

You can either:
  - Provide the scanner name as an argument (recommended), or
  - Leave it blank to be prompted interactively.

The metadata includes supported MIME types, capabilities, vendor information, and more.

Examples:
  # Get metadata for a specific scanner by name
  harbor scanner metadata trivy-scanner

  # Interactively select a scanner if no name is provided
  harbor scanner metadata

Flags:
  --output-format <format>   Output format: 'json' or 'yaml' (default is table view)

```sh
harbor scanner metadata [scanner-name] [flags]
```

### Options

```sh
  -h, --help   help for metadata
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor scanner](harbor-scanner.md)	 - scanner commands

