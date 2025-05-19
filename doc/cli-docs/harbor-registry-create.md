---
title: harbor registry create
weight: 110
---
## harbor registry create

### Description

##### create registry

```sh
harbor registry create [flags]
```

### Examples

```sh
harbor registry create
```

### Options

```sh
      --credential-access-key string      Access key, e.g. user name when credential type is 'basic'
      --credential-access-secret string   Access secret, e.g. password when credential type is 'basic'
      --credential-type string            Credential type, such as 'basic', 'oauth' (default "basic")
      --description string                Description of the registry
  -h, --help                              help for create
      --insecure                          Whether Harbor will verify the server certificate (default true)
      --name string                       Name of the registry
      --type string                       Type of the registry
      --url string                        Registry endpoint URL
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor registry](harbor-registry.md)	 - Manage registries

