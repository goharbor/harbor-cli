---
title: harbor jobservice workers free all
weight: 40
---
## harbor jobservice workers free-all

### Description

##### Free all busy workers (job-id=all)

### Synopsis

Stop all running jobs to free all busy workers.

```sh
harbor jobservice workers free-all [flags]
```

### Examples

```sh
harbor jobservice workers free-all
```

### Options

```sh
  -h, --help   help for free-all
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice workers](harbor-jobservice-workers.md)	 - Manage workers (list all/by pool, free, free-all)

