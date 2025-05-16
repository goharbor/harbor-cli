---
title: harbor quota update
weight: 70
---
## harbor quota update

### Description

##### update quotas for projects

```sh
harbor quota update [QuotaID] [flags]
```

### Options

```sh
  -h, --help                  help for update
      --project-id string     Get quota by project ID
      --project-name string   Get quota by project-name
      --storage string        Enter storage size (e.g., 50GiB, 20MiB, 4TiB)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor quota](harbor-quota.md)	 - Manage quotas

