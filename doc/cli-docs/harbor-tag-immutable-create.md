---
title: harbor tag immutable create
weight: 70
---
## harbor tag immutable create

### Description

##### create immutable tag rule

### Synopsis

create immutable tag rule to the project in harbor

```sh
harbor tag immutable create [flags]
```

### Examples

```sh
harbor tag immutable create
```

### Options

```sh
  -h, --help                     help for create
      --repo-decoration string   repository which either apply or exclude from the rule
      --repo-list string         list of repository to which to either apply or exclude from the rule
      --tag-decoration string    tags which either apply or exclude from the rule
      --tag-list string          list of tags to which to either apply or exclude from the rule
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor tag immutable](harbor-tag-immutable.md)	 - Manage Immutability rules in the project

