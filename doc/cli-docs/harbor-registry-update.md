---
title: harbor registry update
weight: 125
---
## harbor registry update

### Description

##### update registry

```sh
harbor registry update [registry_name] [flags]
```

### Options

```sh
  -k, --credential-access-key string      Access key, e.g. user name when credential type is 'basic'
  -s, --credential-access-secret string   Access secret, e.g. password when credential type is 'basic'
      --credential-type string            Credential type, such as 'basic', 'oauth'
  -d, --description string                Description of the registry
  -h, --help                              help for update
  -i, --insecure                          Whether or not the certificate will be verified when Harbor tries to access the server
  -n, --name string                       Name of the registry
  -t, --type string                       Type of the registry
  -u, --url string                        Registry endpoint URL
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor registry](harbor-registry.md)	 - Manage registries

