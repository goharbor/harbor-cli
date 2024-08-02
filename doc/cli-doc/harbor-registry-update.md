---
title: harbor registry update
weight: 120
---
## harbor registry update

#### update registry

```sh
harbor registry update [flags]
```

### Options

```sh
      --credential-access-key string      Access key, e.g. user name when credential type is 'basic'
      --credential-access-secret string   Access secret, e.g. password when credential type is 'basic'
      --credential-type string            Credential type, such as 'basic', 'oauth'
      --description string                Description of the registry
  -h, --help                              help for update
      --insecure                          Whether or not the certificate will be verified when Harbor tries to access the server (default true)
      --name string                       Name of the registry
      --type string                       Type of the registry
      --url string                        Registry endpoint URL
```

### Options inherited from parent commands

```sh
      --config string          config file (default is $HOME/.harbor/config.yaml) (default "/home/user/.harbor/config.yaml")
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

