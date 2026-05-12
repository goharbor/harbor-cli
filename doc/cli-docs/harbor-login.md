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
  -n, --context-name string   Login context name (optional)
  -h, --help                  help for login
  -p, --password string       Password
      --password-stdin        Take the password from stdin
      --skip-verify-client    Skip whether the clients basic auth credentials shall be validated against the Harbor server during login. This is not recommended as it may lead to storing invalid credentials. Use this flag if you want to skip validation of credentials during login, for example, when the Harbor server is not reachable at the moment of login but you still want to store the credentials for later use.
  -u, --username string       Username
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI

