---
title: harbor scanner update
weight: 90
---
## harbor scanner update

### Description

##### Update a scanner registration

### Synopsis

Update the fields of an existing scanner registration.

You can pass the scanner name as an argument, or the CLI will prompt you to enter a scanner ID.
Only the fields passed through flags will be updated; other fields will retain their existing values.

```sh
harbor scanner update [scanner-name] [flags]
```

### Examples

```sh

  # Update description and URL for a scanner named 'trivy-scanner'
  harbor scanner update trivy-scanner --description "Updated scanner" --url "http://trivy.local:8080"

  # Change the authentication method and credential
  harbor scanner update trivy-scanner --auth Basic --credential "base64encodedAuth"

  # Disable the scanner and rename it
  harbor scanner update trivy-scanner --name "trivy-secure" --disabled

  # If no name is passed, you'll be prompted to enter a Scanner ID
  harbor scanner update --description "Updated via ID prompt"

```

### Options

```sh
      --auth string              Authentication method [None|Basic|Bearer|X-ScannerAdapter-API-Key]
      --credential string        Authorization header for the Scanner Adapter API
      --description string       New description for the scanner
      --disabled                 Disable the scanner registration
  -h, --help                     help for update
      --name string              New name for the scanner
      --skip-cert-verification   Skip certificate verification in HTTP requests
      --url string               Base URL of the scanner adapter
      --use-internal-addr        Use internal registry address for scanning
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor scanner](harbor-scanner.md)	 - scanner commands

