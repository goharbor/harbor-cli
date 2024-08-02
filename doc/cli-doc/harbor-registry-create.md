---
title: harbor registry create
weight: 100
---
## harbor registry create

#### create registry

```sh
harbor registry create [flags]
```

### Options

```sh
      --credential-access-key string      Access key, e.g. user name when credential type is 'basic'
      --credential-access-secret string   Access secret, e.g. password when credential type is 'basic'
      --credential-type string            Credential type, such as 'basic', 'oauth' (default "basic")
      --description string                Description of the registry
  -h, --help                              help for create
      --insecure                          Whether or not the certificate will be verified when Harbor tries to access the server (default true)
      --name string                       Name of the registry
      --type string                       Type of the registry (default "harbor")
      --url string                        Registry endpoint URL
```

### Options inherited from parent commands

```sh
      --config string          config file (default is $HOME/.harbor/config.yaml) (default "/home/user/.harbor/config.yaml")
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

