---
title: harbor apply
weight: 40
---
## harbor apply

### Description

##### Reconcile Harbor with a declarative configuration

### Synopsis

Compare a versioned configuration document with Harbor and apply the
required create and update operations. Apply is additive: resources omitted
from the document are left untouched, and extra live resources are not deleted.

```sh
harbor apply [flags]
```

### Examples

```sh
  harbor apply -f harbor.yaml --dry-run
  harbor apply -f harbor.yaml
  harbor apply -f environments/production
  harbor apply -f harbor.yaml --yes
```

### Options

```sh
      --dry-run       Show the reconciliation plan without changing Harbor
  -f, --file string   YAML/JSON configuration file or directory to apply
  -h, --help          help for apply
  -y, --yes           Apply without an interactive confirmation
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI

