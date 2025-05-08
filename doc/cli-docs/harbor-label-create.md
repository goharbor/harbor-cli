---
title: harbor label create
weight: 80
---
## harbor label create

### Description

##### create label

### Synopsis

create label in harbor

```sh
harbor label create [flags]
```

### Examples

```sh
harbor label create
```

### Options

```sh
      --color string         Color of the label.color is in hex value (default "#FFFFFF")
  -d, --description string   Description of the label
  -h, --help                 help for create
  -n, --name string          Name of the label
  -s, --scope string         Scope of the label. eg- g(global), p(specific project) (default "g")
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor label](harbor-label.md)	 - Manage labels in Harbor

