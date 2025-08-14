---
title: harbor config view
weight: 95
---
## harbor config view

### Description

##### View Harbor configurations

### Synopsis

View Harbor system configurations. You can filter by category:
- authentication: User and service authentication settings
- security: Security policies and certificate settings  
- system: General system behavior and storage settings

```sh
harbor config view [flags]
```

### Options

```sh
      --category string   Filter by category (authentication, security, system)
  -h, --help              help for view
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor config](harbor-config.md)	 - Manage system configurations

