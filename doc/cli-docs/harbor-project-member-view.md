---
title: harbor project member view
weight: 90
---
## harbor project member view

### Description

##### get project member information

```sh
harbor project member view [projectName] [memberID] [flags]
```

### Examples

```sh
  harbor project member view my-project 5
  harbor project member view my-project 5 --wide
```

### Options

```sh
  -h, --help   help for view
      --id     parses projectName as an ID
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project member](harbor-project-member.md)	 - Manage members in a Project

