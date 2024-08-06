---
title: harbor repo delete
weight: 145
---
## harbor repo delete

#### Delete a repository

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
      --config string          config file (default is $HOME/.harbor/config.yaml) (default "/home/user/.harbor/config.yaml")
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

