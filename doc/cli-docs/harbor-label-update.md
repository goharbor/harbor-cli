---
title: harbor label update
weight: 70
---
## harbor label update

### Description

##### update label

```sh
harbor label update [flags]
```

### Examples

```sh
harbor label update [labelname]
```

### Options

```sh
      --color string         Color of the label.color is in hex value
  -d, --description string   Description of the label
      --global               whether to list global or project scope labels. (default scope is global)
  -h, --help                 help for update
  -n, --name string          Name of the label
  -p, --project string       project name when query project labels
  -i, --project-id int       project ID when query project labels
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor label](harbor-label.md)	 - Manage labels in Harbor

