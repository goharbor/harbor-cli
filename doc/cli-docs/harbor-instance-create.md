---
title: harbor instance create
weight: 50
---
## harbor instance create

### Description

##### create instance

```sh
harbor instance create [flags]
```

### Options

```sh
  -a, --authmode string      Choosing different types of authentication method (default "NONE")
      --description string   Description of the instance
      --enable               Whether it is enable or not (default true)
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

* [harbor instance](harbor-instance.md)	 - Manage instance in Harbor

