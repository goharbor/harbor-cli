---
title: harbor logs events
weight: 90
---
## harbor logs events

### Description

##### List supported Harbor audit log event types

### Synopsis

List supported Harbor audit log event types.

By default, all event types are shown.
Use --page and --page-size to paginate the result.

Examples:
  harbor-cli logs events
  harbor-cli logs events --page 2 --page-size 5
  harbor-cli logs events --output-format json --page 2 --page-size 5

```sh
harbor logs events [flags]
```

### Options

```sh
  -h, --help            help for events
      --page int        Page number (default 1)
      --page-size int   Size of per page (default 10)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor logs](harbor-logs.md)	 - Get recent logs of the projects which the user is a member of

