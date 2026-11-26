---
title: harbor auditlog list
weight: 10
---
## harbor auditlog list

### Description

##### List audit logs

### Synopsis

List system-wide audit logs in Harbor

```sh
harbor auditlog list [flags]
```

### Examples

```sh
harbor auditlog list --page 1 --page-size 10 --username admin
```

### Options

```sh
  -h, --help               help for list
      --operation string   Filter logs by operation type (e.g. create, delete, pull)
      --page int           Page number (default 1)
      --page-size int      Size of per page (default 10)
  -q, --query string       Query string for filtering logs
  -r, --resource string    Filter logs by target resource
  -u, --username string    Filter logs by username
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor auditlog](harbor-auditlog.md)	 - Manage and view audit logs

