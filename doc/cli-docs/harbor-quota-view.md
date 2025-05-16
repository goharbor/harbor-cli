---
title: harbor quota view
weight: 30
---
## harbor quota view

### Description

##### get quota by quota ID

```sh
harbor quota view [quotaID] [flags]
```

### Options

```sh
  -h, --help                  help for view
      --project-id string     Get quota by project ID
      --project-name string   Get quota by project-name
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor quota](harbor-quota.md)	 - Manage quotas

