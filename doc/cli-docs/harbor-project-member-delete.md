---
title: harbor project member delete
weight: 30
---
## harbor project member delete

### Description

##### delete member by username

### Synopsis

delete members in a project by username of the member

```sh
harbor project member delete [flags]
```

### Examples

```sh
  harbor project member delete my-project --username user
```

### Options

```sh
  -a, --all               Deletes all members of the project
  -h, --help              help for delete
      --id                parses projectName as an ID
  -u, --username string   Username of the member
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project member](harbor-project-member.md)	 - Manage members in a Project

