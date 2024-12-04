---
title: harbor login
weight: 15
---
## harbor login

### Description

##### Log in to Harbor registry

Authenticate with Harbor Registry. Depending on how the command is invoked, it behaves differently:

##### With `-u` / `-p` Flags and `server`
  - Opens the login view to obtain new credentials.
  - Updates the config file with the new credentials.
  - If the specified credential name already exists, it updates the existing entry.
  - If the credential name does not exist, it adds a new entry for the credential.

##### Without `-u` / `-p` Flags and `server`
a. No Existing Credentials in Config:
  - Opens the login view to input credentials.
  - Stores the entered credentials in the config file.

b. Existing Credentials in Config:
  - Uses the stored credentials from the config file.
  - Skips the login view and proceeds to authenticate using the existing credentials.

For more info on the harbor-cli config management see the [harbor config docs](harbor-config.md)

### Synopsis

Authenticate with Harbor Registry.

```sh
harbor login [server] [flags]
```

### Options

```sh
  -h, --help              help for login
      --name string       name for the set of credentials
  -p, --password string   Password
  -u, --username string   Username
```

### Options inherited from parent commands

```sh
      --config string          config file (default is $HOME/.harbor/config.yaml) (default "/home/user/.harbor/config.yaml")
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI
* [harbor config](harbor-config.md) - Harbor Config Management

