---
title: harbor robot list
weight: 65
---
## harbor robot list

### Description

##### list robot

### Synopsis

List robot accounts in Harbor.

This command displays a list of system-level robot accounts. The list includes basic
information about each robot account, such as ID, name, creation time, and
expiration status.

System-level robots have permissions that can span across multiple projects, making
them suitable for CI/CD pipelines and automation tasks that require access to 
multiple projects in Harbor.

You can control the output using pagination flags and format options:
- Use --page and --page-size to navigate through results
- Use --sort to order the results by name, creation time, etc.
- Use -q/--query to filter robots by specific criteria
- Set output-format in your configuration for JSON, YAML, or other formats

Examples:
  # List all system robots
  harbor-cli robot list

  # List system robots with pagination
  harbor-cli robot list --page 2 --page-size 20

  # List system robots with custom sorting
  harbor-cli robot list --sort name

  # Filter system robots by name
  harbor-cli robot list -q name=ci-robot

  # Get robot details in JSON format
  harbor-cli robot list --output-format json

```sh
harbor robot list [projectName] [flags]
```

### Options

```sh
  -h, --help            help for list
      --page int        Page number (default 1)
      --page-size int   Size of per page (default 10)
  -q, --query string    Query string to query resources
      --sort string     Sort the resource list in ascending or descending order
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor robot](harbor-robot.md)	 - Manage robot accounts

