---
title: harbor scanner create
weight: 95
---
## harbor scanner create

### Description

##### Create a scanner

```sh
harbor scanner create [flags]
```

### Options

```sh
      --auth string                           Authentication method [None|Basic|Bearer|X-ScannerAdapter-API-Key]
      --credential string                     Authorization header for the Scanner Adapter API
      --description string                    New description for the scanner
      --disabled[=true|false]                 Disable the scanner registration
  -h, --help[=true|false]                     help for create
      --name string                           New name for the scanner
      --ping[=true|false]                     Ping the scanner adapter without creating it.
      --skip-cert-verification[=true|false]   Skip certificate verification in HTTP requests
      --url string                            Base URL of the scanner adapter
      --use-internal-addr[=true|false]        Use internal registry address for scanning
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -l, --log-format string      Output format for logging. One of: json|text (default "text")
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose[=true|false]   verbose output
```

### SEE ALSO

* [harbor scanner](harbor-scanner.md)	 - scanner commands

