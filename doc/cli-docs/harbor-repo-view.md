---
title: harbor repo view
weight: 55
---
## harbor repo view

### Description

##### Get repository information

### Synopsis

Get information of a particular repository in a project

```sh
harbor repo view [flags]
```

### Examples

```sh
  harbor repo view <project_name>/<repo_name>
```

### Options

```sh
  -h, --help   help for view
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor repo](harbor-repo.md)	 - Manage repositories

