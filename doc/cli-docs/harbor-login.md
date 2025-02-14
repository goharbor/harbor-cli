---
title: harbor login
weight: 15
---
## harbor login

### Description

##### Log in to Harbor registry

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

