---
title: harbor quota
weight: 30
---
## harbor quota

### Description

##### Manage quotas

### Synopsis

Manage quotas of projects

### Examples

```sh
  harbor quota list
```

### Options

```sh
  -h, --help   help for quota
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -l, --log-format string      Output format for logging. One of: json|text (default "text")
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI
* [harbor quota list](harbor-quota-list.md)	 - list quotas
* [harbor quota update](harbor-quota-update.md)	 - update quotas for projects
* [harbor quota view](harbor-quota-view.md)	 - get quota by quota ID

