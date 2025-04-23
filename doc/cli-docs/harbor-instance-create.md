---
title: harbor instance create
weight: 0
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
  harbor-cli instance create --name my-instance --provider Dragonfly --url http://dragonfly.local --description "My preheat provider instance" --enable=true
```

### Options

```sh
  -a, --authmode string      Choosing different types of authentication method (default "NONE")
      --description string   Description of the instance
      --enable               Whether it is enabled or not (default true)
  -h, --help                 help for create
  -i, --insecure             Whether or not the certificate will be verified when Harbor tries to access the server (default true)
  -n, --name string          Name of the instance
  -p, --provider string      Provider for the instance
  -u, --url string           URL for the instance
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor instance](harbor-instance.md)	 - Manage preheat provider instances in Harbor

