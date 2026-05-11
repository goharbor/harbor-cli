---
title: harbor project summary
weight: 15
---
## harbor project summary

### Description

##### Get summary of a project

### Synopsis

Get summary of a project by name or ID. If no arguments are provided, it will prompt for the project name. Use --id to specify the project ID instead of the name.

```sh
harbor project summary [NAME|ID] [flags]
```

### Examples

```sh
harbor project summary my-project or harbor project summary 1 --id
```

### Options

```sh
  -h, --help   help for summary
      --id     Get project by id
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project](harbor-project.md)	 - Manage projects and assign resources to them

