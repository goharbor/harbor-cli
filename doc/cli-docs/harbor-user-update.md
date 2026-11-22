---
title: harbor user update
weight: 5
---
## harbor user update

### Description

##### update user profile

### Synopsis

The 'update' command allows sysadmins to modify an existing user's profile information, such as their email, realname, or comment.

```sh
harbor user update [USER_NAME_OR_ID] [flags]
```

### Examples

```sh
  harbor user update admin --email newadmin@example.com --realname "System Admin"
```

### Options

```sh
      --comment string    Comment
      --email string      Email
  -h, --help              help for update
      --realname string   Realname
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor user](harbor-user.md)	 - Manage users

