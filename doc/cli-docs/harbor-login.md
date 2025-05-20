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
  -p, --password string   Password
      --password-stdin    Take the password from stdin
  -u, --username string   Username
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI

