---
title: harbor label delete
weight: 30
---
## harbor label delete

### Description

##### delete label

```sh
harbor label delete [flags]
```

### Examples

```sh
harbor label delete [labelname]
```

### Options

```sh
      --global           whether to list global or project scope labels. (default scope is global)
  -h, --help             help for delete
  -p, --project string   project name when query project labels
  -i, --project-id int   project ID when query project labels
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor label](harbor-label.md)	 - Manage labels in Harbor

