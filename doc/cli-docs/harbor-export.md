---
title: harbor export
weight: 5
---
## harbor export

### Description

##### Export Harbor API-managed configuration

### Synopsis

Export Harbor API-managed configuration as a versioned declarative document.

The export contains system settings, registry endpoints, projects, quotas,
webhooks, and replication policies. It does not contain credentials, robot
secrets, users, artifacts, database contents, or Harbor deployment settings.

```sh
harbor export [flags]
```

### Examples

```sh
  harbor export -f harbor.yaml
  harbor export -f harbor.json -o json
  harbor export -o yaml > harbor.yaml
```

### Options

```sh
  -f, --file string   Write configuration to a YAML or JSON file (default: stdout)
  -h, --help          help for export
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI

