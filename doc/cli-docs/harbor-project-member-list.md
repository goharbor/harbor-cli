---
title: harbor project member list
weight: 35
---
## harbor project member list

### Description

##### list members in a project

### Synopsis

list members in a project by projectName

```sh
harbor project member list [projectName] [flags]
```

### Examples

```sh
  harbor project member list my-project
```

### Options

```sh
  -h, --help            help for list
  -n, --name string     Member Name to search
      --page int        Page number (default 1)
      --page-size int   Size of per page (default 10)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project member](harbor-project-member.md)	 - Manage members in a Project

