---
title: harbor tag retention tasks list
weight: 50
---
## harbor tag retention tasks list

### Description

##### List retention tasks for an execution

### Synopsis

List repository-level retention tasks for a specific retention execution

```sh
harbor tag retention tasks list [PROJECT_NAME] [flags]
```

### Options

```sh
  -e, --execution-id int      Retention execution ID
  -h, --help                  help for list
  -i, --project-id int        Project ID
  -p, --project-name string   Project name
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -l, --log-format string      Output format for logging. One of: json|text (default "text")
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor tag retention tasks](harbor-tag-retention-tasks.md)	 - Manage retention execution tasks

