---
title: harbor instance create
weight: 10
---
## harbor instance create

### Description

##### Create a new preheat provider instance in Harbor

### Synopsis

Create a new preheat provider instance within Harbor for distributing container images.
The instance can be an external service such as Dragonfly, Kraken, or any custom provider.
You will need to provide the instance's name, vendor, endpoint, and optionally other details such as authentication and security options.

```sh
harbor instance create [flags]
```

### Examples

```sh
  harbor-cli instance create --name my-instance --provider dragonfly --url http://dragonfly.local --description "My preheat provider instance" --enable=true
```

### Options

```sh
      --auth-password string   Password for BASIC authentication
      --auth-token string      Token for OAUTH authentication
      --auth-username string   Username for BASIC authentication
  -a, --authmode string        Authentication mode (NONE, BASIC, OAUTH) (default "NONE")
  -d, --description string     Description of the instance
      --enable                 Whether the instance is enabled or not (default true)
  -h, --help                   help for create
  -i, --insecure               Whether or not the certificate will be verified when Harbor tries to access the server
  -n, --name string            Name of the instance
  -p, --provider string        Provider for the instance (e.g. dragonfly, kraken)
  -u, --url string             Endpoint URL for the instance
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor instance](harbor-instance.md)	 - Manage preheat provider instances in Harbor

