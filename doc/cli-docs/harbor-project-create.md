---
title: harbor project create
weight: 80
---
## harbor project create

### Description

##### create project

```sh
harbor project create [project name] [flags]
```

### Options

```sh
  -h, --help                   help for create
      --proxy-cache            Whether the project is a proxy cache project
      --public                 Project is public or private
      --registry-id string     ID of referenced registry when creating the proxy cache project
      --storage-limit string   Storage quota of the project
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project](harbor-project.md)	 - Manage projects and assign resources to them

