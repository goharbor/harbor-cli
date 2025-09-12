---
title: harbor project member create
weight: 15
---
## harbor project member create

### Description

##### create project member

### Synopsis

create project member by Name

```sh
harbor project member create [flags]
```

### Examples

```sh
  harbor project member create my-project --username user --role Developer
```

### Options

```sh
      --groupid int        Group ID
      --groupname string   Group Name
      --grouptype int      Group Type
  -h, --help               help for create
      --id                 parses projectName as an ID
      --ldapdn string      DN of LDAP Group
      --role string        Role Name [one of Project_Admin, Developer, Guest, Maintainer, Limited_Guest]
      --roleid int         Role ID
      --userid int         User ID
      --username string    Username
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project member](harbor-project-member.md)	 - Manage members in a Project

