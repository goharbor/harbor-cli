---
title: harbor tag retention create
weight: 65
---
## harbor tag retention create

### Description

##### create retention policy

### Synopsis

create a retention policy for a project

```sh
harbor tag retention create [PROJECT_NAME] [flags]
```

### Options

```sh
      --cron string              schedule cron expression for retention policy
      --dry-run                  validate policy file and project scope without creating
  -f, --file string              retention policy file in JSON format (optional in interactive mode)
  -h, --help                     help for create
      --keep-latest int          number of most recently pushed artifacts to retain
      --project string           project name
      --repo-decoration string   repository selector decoration: repoMatches or repoExcludes
      --repo-list string         repository selector pattern, for example **
      --tag-decoration string    tag selector decoration: matches or excludes
      --tag-list string          tag selector pattern, for example **
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor tag retention](harbor-tag-retention.md)	 - Manage retention policies

