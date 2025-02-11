---
title: harbor user reset
weight: 10
---
## harbor user reset

### Description

##### reset user's password

### Synopsis

Resets the password for a specific user by providing their username

```sh
harbor user reset [username] [flags]
```

### Options

```sh
  -h, --help         help for reset
      --userID int   ID of the user (default -1)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor user](harbor-user.md)	 - Manage users

