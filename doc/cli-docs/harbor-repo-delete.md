---
title: harbor repo delete
weight: 160
---
## harbor repo delete

### Description

##### Delete a repository

### Synopsis

Delete a repository within a project in Harbor

```sh
harbor repo delete [flags]
```

### Examples

```sh
  harbor repository delete [project_name]/[repository_name]
```

### Options

```sh
  -h, --help   help for delete
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -l, --log-format string      Output format for logging. One of: json|text (default "text")
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor repo](harbor-repo.md)	 - Manage repositories

