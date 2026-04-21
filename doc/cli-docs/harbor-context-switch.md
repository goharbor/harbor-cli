---
title: harbor context switch
weight: 50
---
## harbor context switch

### Description

##### Switch to a new context

```sh
harbor context switch <none|context> [flags]
```

### Examples

```sh
harbor context switch harbor-cli@https-demo-goharbor-io
```

### Options

```sh
  -h, --help   help for switch
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -l, --log-format string      Output format for logging. One of: json|text (default "text")
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor context](harbor-context.md)	 - Manage locally available contexts

