---
title: harbor user delete
weight: 185
---
## harbor user delete

### Description

##### delete user by name or id

### Synopsis

Delete one or more users by username or user ID, or pick one interactively when no argument is given.

Each argument is matched against usernames first. A number that is not a username is used as a user ID, and a number that is the username of one user and the ID of another is rejected as ambiguous. Every argument is resolved before any user is deleted, so nothing is deleted if one of them does not resolve.

```sh
harbor user delete [username|id]... [flags]
```

### Options

```sh
  -h, --help[=true|false]   help for delete
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -l, --log-format string      Output format for logging. One of: json|text (default "text")
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose[=true|false]   verbose output
```

### SEE ALSO

* [harbor user](harbor-user.md)	 - Manage users

