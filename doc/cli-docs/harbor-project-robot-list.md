---
title: harbor project robot list
weight: 35
---
## harbor project robot list

### Description

##### list robot

### Synopsis

List robot accounts in a Harbor project.

This command displays a list of robot accounts, either from a specific project
or by prompting you to select a project interactively. The list includes basic
information about each robot account, such as ID, name, creation time, and
expiration status.

The command supports multiple ways to specify the project:
- By providing a project name as an argument
- By using the --project-id flag
- By using the -q/--query flag with a project filter
- Without any arguments, which will prompt for project selection

You can control the output using pagination flags and format options:
- Use --page and --page-size to navigate through results
- Use --sort to order the results
- Set output-format in your configuration for JSON, YAML, or other formats

Examples:
  # List robots in a specific project by name
  harbor-cli project robot list myproject

  # List robots in a project by ID
  harbor-cli project robot list --project-id 123

  # List robots with pagination
  harbor-cli project robot list --page 2 --page-size 20

  # List robots with custom sorting
  harbor-cli project robot list --sort name

  # Interactive listing (will prompt for project selection)
  harbor-cli project robot list

```sh
harbor project robot list [projectName] [flags]
```

### Options

```sh
  -h, --help             help for list
      --page int         Page number (default 1)
      --page-size int    Size of per page (default 10)
      --project-id int   Project ID
  -q, --query string     Query string to query resources
      --sort string      Sort the resource list in ascending or descending order
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project robot](harbor-project-robot.md)	 - Manage robot accounts

