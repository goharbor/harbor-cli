---
title: harbor scanner view
weight: 95
---
## harbor scanner view

### Description

##### Display detailed information about a scanner registration

### Synopsis

Display full details of a scanner registration in Harbor.

You can:
  - Provide a scanner name to view its details directly.
  - Omit the argument to select a scanner interactively by ID.

Supports custom output formats using the --output-format flag (e.g., json, yaml, table).

Examples:
  # View a specific scanner by name
  harbor scanner view trivy-scanner

  # Interactively choose a scanner to view
  harbor scanner view

  # View scanner in JSON format
  harbor scanner view trivy-scanner --output-format=json

```sh
harbor scanner view [scanner-name] [flags]
```

### Options

```sh
  -h, --help   help for view
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor scanner](harbor-scanner.md)	 - scanner commands

