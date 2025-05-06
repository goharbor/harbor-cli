---
title: harbor quota list
weight: 85
---
## harbor quota list

### Description

##### list quotas

### Synopsis

list quotas specified for each project

```sh
harbor quota list [flags]
```

### Options

```sh
  -h, --help            help for list
      --page int        Page number (default 1)
      --page-size int   Size of per page (use 0 to fetch all)
      --ref string      Reference type of quota
      --refid string    Reference ID of quota
      --sort string     Sort the resource list in ascending or descending order
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor quota](harbor-quota.md)	 - Manage quotas

