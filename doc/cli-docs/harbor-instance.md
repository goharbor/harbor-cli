---
title: harbor instance
weight: 5
---
## harbor instance

### Description

##### Manage preheat provider instances in Harbor

### Synopsis

Manage preheat provider instances used by Harbor for pre-distributing container images.
These instances represent external services such as Dragonfly or Kraken that help preheat images across nodes.

### Options

```sh
  -h, --help   help for instance
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI
* [harbor instance create](harbor-instance-create.md)	 - Create a new preheat provider instance in Harbor
* [harbor instance delete](harbor-instance-delete.md)	 - Delete a preheat provider instance by its name or ID
* [harbor instance list](harbor-instance-list.md)	 - List all preheat provider instances in Harbor

