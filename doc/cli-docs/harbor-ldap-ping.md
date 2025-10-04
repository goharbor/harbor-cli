---
title: harbor ldap ping
weight: 90
---
## harbor ldap ping

### Description

##### ping ldap server

```sh
harbor ldap ping [flags]
```

### Options

```sh
  -h, --help                    help for ping
      --ldap-base-dn string     The base dn from which to lookup the user
      --ldap-filter string      Search Filter of ldap service
      --ldap-password string    search password of the ldap service
      --ldap-scope int          search scope of ldap service default 0 base, 1 OneLevel, 2 Subtree.
      --ldap-search-dn string   User's dn who has the permission to search the ldap server
      --ldap-uid string         attribute used in search to match the user. It could be cn, uid based on your LDAP/AD.
      --ldap-url string         URL of the ldap service
      --ldap-verify             Verify Ldap server certificate
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor ldap](harbor-ldap.md)	 - Manage ldap users and groups

