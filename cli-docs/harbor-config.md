---
title: harbor config
weight: 5
---
## harbor config

### Description

##### Manage system configurations

### Synopsis

Manage Harbor system configurations including viewing, exporting, and applying settings.

Configuration management workflow:
1. View configurations in table format or export to files
2. Edit exported configuration files as needed  
3. Apply modified configurations back to Harbor

Categories available:
- authentication (auth): LDAP, OIDC, UAA authentication settings
- security (sec): Security policies, certificates, and access control
- system (sys): General system behavior, storage, and operational settings

### Examples

```sh
  # View configurations
  harbor config view                          # Table view of all configs
  harbor config view -c auth                  # View only authentication configs
  harbor config view -c sec                   # View only security configs

  # Export configurations to files
  harbor config view -o json > config.json                    # Export all configs as JSON
  harbor config view -c auth -o yaml | tee auth-config.yaml   # Export auth configs as YAML
  harbor config view -c sys -o json > system-config.json     # Export system configs as JSON

  # Apply configurations from files
  harbor config apply -f config.json         # Apply complete configuration
  harbor config apply -f auth-config.yaml    # Apply only authentication settings
  
  # Configuration backup and restore workflow  
  harbor config view -o yaml > backup.yaml   # Create backup
  # ... make changes to Harbor via UI or other means ...
  harbor config apply -f backup.yaml         # Restore from backup
```

### Options

```sh
  -h, --help   help for config
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI
* [harbor config apply](harbor-config-apply.md)	 - Update system configurations from local config file
* [harbor config view](harbor-config-view.md)	 - View Harbor configurations

