---
title: harbor repo info
weight: 140
---
## harbor repo info

#### Get repository information

### Synopsis

Get information of a particular repository in a project

```sh
harbor repo info [flags]
```

### Examples

```sh
  harbor repo info <project_name>/<repo_name>
```

### Options

```sh
  -h, --help   help for info
```

### Options inherited from parent commands

```sh
      --config string          config file (default is $HOME/.harbor/config.yaml) (default "/home/user/.harbor/config.yaml")
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

