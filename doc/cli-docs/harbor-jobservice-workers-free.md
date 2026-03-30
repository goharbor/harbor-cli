---
title: harbor jobservice workers free
weight: 20
---
## harbor jobservice workers free

### Description

##### Free one worker (--job-id required)

### Synopsis

Stop a running job by job ID to free its worker.

```sh
harbor jobservice workers free [flags]
```

### Examples

```sh
harbor jobservice workers free --job-id abc123
```

### Options

```sh
  -h, --help            help for free
      --job-id string   Running job ID to stop
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor jobservice workers](harbor-jobservice-workers.md)	 - Manage workers (list all/by pool, free, free-all)

