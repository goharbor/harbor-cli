---
title: harbor user update
weight: 55
---
## harbor user update

### Description

##### update user's profile

### Synopsis

Update a user's profile by providing their username, allowing you to modify personal details such as realname, email, or comment

```sh
harbor user update [flags]
```

### Examples

```sh
harbor user update [username]
```

### Options

```sh
      --comment string    Comment of the user
      --email string      Email of the user
  -h, --help              help for update
      --realname string   Realname of the user
      --userID int        ID of the user (default -1)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor user](harbor-user.md)	 - Manage users

