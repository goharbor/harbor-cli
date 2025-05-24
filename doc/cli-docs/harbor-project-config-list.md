---
title: harbor project config list
weight: 85
---
## harbor project config list

### Description

##### List configuration of a Harbor project by name or ID

### Synopsis

Display the configuration metadata of a Harbor project specified by its name or ID.

If no project name or ID is provided as an argument, you will be prompted to select a project interactively.

You can use the global flag '--output-format' to specify the output format, e.g. 'json' or 'yaml', for machine-readable output.

Examples:

  # List configuration of project 'myproject' by name
  harbor-cli project config list myproject

  # List configuration of project with ID '123'
  harbor-cli project config list 123

  # Run interactively (prompt to select project)
  harbor-cli project config list

  # List config in JSON format
  harbor-cli project config list myproject --output-format json


```sh
harbor project config list  [project_name] [flags]
```

### Options

```sh
  -h, --help   help for list
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
      --id                     Use project ID instead of name
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project config](harbor-project-config.md)	 - Manage project configuration

