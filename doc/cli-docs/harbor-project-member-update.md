---
title: harbor project member update
weight: 20
---
## harbor project member update

### Description

##### update member by ID or name

### Synopsis

update member in a project by MemberID

```sh
harbor project member update [ProjectName Or ID] [member ID] [flags]
```

### Examples

```sh
  harbor project member update my-project [memberID] --roleid 2
```

### Options

```sh
  -h, --help         help for update
      --id int       Member ID
      --roleid int   Role to be updated
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project member](harbor-project-member.md)	 - Manage members in a Project

