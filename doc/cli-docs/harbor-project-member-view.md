---
title: harbor project member view
weight: 90
---
## harbor project member view

### Description

##### get project member by ID or name

### Synopsis

get member details by MemberID

```sh
harbor project member view [ProjectName Or ID] [member ID] [flags]
```

### Examples

```sh
  harbor project member view my-project [memberID]
```

### Options

```sh
  -h, --help                 help for view
      --id int               Member ID
  -p, --projectname string   Project Name
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project member](harbor-project-member.md)	 - Manage members in a Project

