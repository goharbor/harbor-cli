---
title: harbor scanner set default
weight: 25
---
## harbor scanner set-default

### Description

##### Set the default scanner for Harbor

### Synopsis

Set the scanner that will be used as the default in Harbor. This scanner will be used for all default scanning tasks unless another scanner is specified.

```sh
harbor scanner set-default [flags]
```

### Examples

```sh
harbor scanner set-default [scanner-name]
		OR 
		harbor scanner set-default --id <scanner-id>
```

### Options

```sh
  -h, --help   help for set-default
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor scanner](harbor-scanner.md)	 - scanner commands

