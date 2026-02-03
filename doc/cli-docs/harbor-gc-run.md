---
title: harbor gc run
weight: 55
---
## harbor gc run

### Description

##### Run Garbage Collection manually

```sh
harbor gc run [flags]
```

### Options

```sh
      --delete-untagged   Delete untagged artifacts (default true)
      --dry-run           Simulate GC without deleting artifacts
  -h, --help              help for run
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor gc](harbor-gc.md)	 - Manage Garbage Collection

