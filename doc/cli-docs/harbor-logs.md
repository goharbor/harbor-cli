---
title: harbor logs
weight: 30
---
## harbor logs

### Description

##### Get recent logs of the projects which the user is a member of

### Synopsis

Get recent logs of the projects which the user is a member of.
This command retrieves the audit logs for the projects the user is a member of. It supports pagination, sorting, and filtering through query parameters. The logs can be followed in real-time with the --follow flag, and the output can be formatted as JSON with the --output-format flag.

When --page and/or --page-size are explicitly provided, a pagination summary (for example: "Showing 6-10 of 14") is shown in default table output.

Convenience filter flags are available to build query expressions:
- --operation (operation type, e.g., create, delete, pull, login, logout, update)
- --resource-type (resource type, e.g., user, artifact, project, repository)
- --resource (resource name)
- --username (username)
- --from-time and optional --to-time (for op_time range)

IMPORTANT: Event types from 'harbor logs events' are compound strings (e.g., 'create_user', 'delete_artifact').
These must be split into separate flags:
  create_user       → --operation create --resource-type user
  delete_artifact   → --operation delete --resource-type artifact
  pull_repository   → --operation pull --resource-type repository
  login_user        → --operation login --resource-type user

Examples:
harbor-cli logs --operation create --resource-type user
harbor-cli logs --operation delete --resource-type artifact --resource library/nginx
harbor-cli logs --page 1 --page-size 10 --query "operation=push" --sort "op_time:desc"
harbor-cli logs --follow --refresh-interval 2s
harbor-cli logs --output-format json

```sh
harbor logs [flags]
```

### Options

```sh
  -f, --follow                    Follow log output (tail -f behavior)
      --from-time string          Start timestamp for op_time range (RFC3339 or 'YYYY-MM-DD HH:MM:SS'). Required when using --to-time
  -h, --help                      help for logs
      --operation string          Filter by operation (e.g. create, delete, pull, login, logout, update). Event types from 'harbor logs events' must be split: use --operation and --resource-type separately.
      --page int                  Page number (default 1)
      --page-size int             Size of per page (default 10)
  -q, --query string              Query string to query resources
  -n, --refresh-interval string   Interval to refresh logs when following (default: 5s)
      --resource string           Filter by resource name
      --resource-type string      Filter by resource type (e.g. user, artifact, project, repository). Event types from 'harbor logs events' must be split: use --operation and --resource-type separately.
      --sort string               Sort the resource list in ascending or descending order
      --to-time string            End timestamp for op_time range (RFC3339 or 'YYYY-MM-DD HH:MM:SS'). Optional when --from-time is set; defaults to current time
      --username string           Filter by username
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI
* [harbor logs events](harbor-logs-events.md)	 - List supported Harbor audit log event types

