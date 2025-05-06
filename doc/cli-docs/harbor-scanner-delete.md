---
title: harbor scanner delete
weight: 0
---
## harbor scanner delete

### Description

##### Delete a scanner registration

### Synopsis

Delete a scanner registration from Harbor.

You can:
  - Provide the scanner name as an argument to delete it directly, or
  - Omit the argument to select a scanner interactively.

Note: Deleting a scanner will permanently remove its registration and associated metadata from the system.

Examples:
  # Delete a scanner by name
  harbor scanner delete trivy-scanner

  # Interactively choose a scanner to delete
  harbor scanner delete

```sh
harbor scanner delete [scanner-name] [flags]
```

### Options

```sh
  -h, --help   help for delete
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor scanner](harbor-scanner.md)	 - scanner commands

