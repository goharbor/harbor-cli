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

NOTE: Each event type shown maps to two separate filter flags:
  create_user        → --operation create --resource-type user
  create_artifact    → --operation create --resource-type artifact
  delete_artifact    → --operation delete --resource-type artifact
  delete_project     → --operation delete --resource-type project
  login_user         → --operation login --resource-type user
  logout_user        → --operation logout --resource-type user
  pull_repository    → --operation pull --resource-type repository
  push_artifact      → --operation push --resource-type artifact
  update_artifact    → --operation update --resource-type artifact

When filtering logs with 'harbor logs', use both --operation and --resource-type flags separately.

Examples:
  harbor-cli logs events
  harbor-cli logs events --page 2 --page-size 5
  harbor-cli logs events --output-format json --page 2 --page-size 5
  harbor-cli logs --operation create --resource-type user
  harbor-cli logs --operation delete --resource-type artifact

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
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor logs](harbor-logs.md)	 - Get recent logs of the projects which the user is a member of

