---
title: harbor project create
weight: 80
---
## harbor project create

### Description

##### create project

```sh
harbor project create [flags]
```

### Options

```sh
  -h, --help                   help for create
      --name string            Name of the project
      --proxy-cache            Whether the project is a proxy cache project
      --public                 Project is public or private (default true)
      --registry-id string     ID of referenced registry when creating the proxy cache project
      --storage-limit string   Storage quota of the project (default "-1")
```

### Options inherited from parent commands

```sh
      --config string          config file (default is $HOME/.harbor/config.yaml) (default "/home/user/.harbor/config.yaml")
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project](harbor-project.md)	 - Manage projects and assign resources to them

