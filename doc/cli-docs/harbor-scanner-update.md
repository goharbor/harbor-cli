---
title: harbor scanner update
weight: 95
---
## harbor scanner update

### Description

##### update scanner

```sh
harbor scanner update [flags]
```

### Options

```sh
      --auth string   Authentication approach of the scanner [None|Basic|Bearer|X-ScannerAdapter-API-Key]
      --cred string   HTTP Authorization header value sent with each request to the Scanner Adapter API
      --des string    Description of the scanner
      --disable       Indicate whether the registration is enabled or not
  -h, --help          help for update
      --internal      Indicate whether use internal registry addr for the scanner to pull content or not
      --name string   Name of the scanner
      --skip          Indicate if skip the certificate verification when sending HTTP requests
      --url string    Base URL of the scanner adapter
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor scanner](harbor-scanner.md)	 - scanner commands

