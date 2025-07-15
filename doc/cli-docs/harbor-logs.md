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

harbor-cli logs --page 1 --page-size 10 --query "operation=push" --sort "op_time:desc"

harbor-cli logs --follow --refresh-interval 2s

harbor-cli logs --output-format json

```sh
harbor logs [flags]
```

### Options

```sh
  -f, --follow                    Follow log output (tail -f behavior)
  -h, --help                      help for logs
      --page int                  Page number (default 1)
      --page-size int             Size of per page (default 10)
  -q, --query string              Query string to query resources
  -n, --refresh-interval string   Interval to refresh logs when following (default: 5s)
      --sort string               Sort the resource list in ascending or descending order
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI

