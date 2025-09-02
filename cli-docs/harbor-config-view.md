---
title: harbor config view
weight: 60
---
## harbor config view

### Description

##### View Harbor configurations

### Synopsis

View Harbor system configurations. You can filter by category using full names or shorthand:

Categories:
- authentication (auth): User and service authentication settings (LDAP, OIDC, UAA)
- security (sec): Security policies and certificate settings
- system (sys): General system behavior and storage settings

Examples:
  harbor config view                        # View all configurations
  harbor config view --category auth        # View authentication configs
  harbor config view --cat sec              # View security configs (shorthand)
  harbor config view --cat sys              # View system configs

  # Export configurations to files
  harbor config view -o json > config.json                    # Save all configs as JSON
  harbor config view --cat auth -o yaml | tee auth-config.yaml   # Save auth configs as YAML and display
  harbor config view --cat sec -o json > security-config.json   # Save security configs as JSON

```sh
harbor config view [flags]
```

### Options

```sh
      --cat string        Filter by category (shorthand for --category)
      --category string   Filter by category: authentication (auth), security (sec), system (sys)
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

