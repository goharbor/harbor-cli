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
  -f, --fuzzy string    Fuzzy search for member with name
  -h, --help            help for list
      --id              Parses projectName as an ID
      --page int        Page number (default 1)
      --page-size int   Size of per page (default 10)
  -s, --search string   Search for member with name
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project member](harbor-project-member.md)	 - Manage members in a Project

