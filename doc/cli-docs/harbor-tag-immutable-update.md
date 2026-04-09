---
title: harbor tag immutable update
weight: 50
---
## harbor tag immutable update

### Description

##### update immutable tag rule

### Synopsis

update immutable tag rule for a project in harbor

```sh
harbor tag immutable update [PROJECT_NAME] [flags]
```

### Options

```sh
  -h, --help                     help for update
      --immutable-id int         immutable rule ID to update
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

