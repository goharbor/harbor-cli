---
title: harbor logout
weight: 60
---
## harbor logout

### Description

##### Log out from Harbor registry

### Synopsis

Remove the current credential from the local CLI config.

```sh
harbor logout [flags]
```

### Examples

```sh
  harbor logout
```

### Options

```sh
  -h, --help   help for logout
  -y, --yes    Skip confirmation prompt
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI

