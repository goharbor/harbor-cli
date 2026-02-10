---
title: harbor usergroup create
weight: 30
---
## harbor usergroup create

### Description

##### create user group

```sh
harbor usergroup create [flags]
```

### Options

```sh
  -h, --help             help for create
  -l, --ldap-dn string   The DN of the LDAP group if group type is 1 (LDAP group)
  -n, --name string      Group name
  -t, --type int         Group type
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor usergroup](harbor-usergroup.md)	 - Manage usergroup

