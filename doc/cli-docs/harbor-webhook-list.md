---
title: harbor webhook list
weight: 60
---
## harbor webhook list

### Description

##### List all webhook policies for a Harbor project

### Synopsis

This command retrieves and displays all webhook policies associated with a Harbor project.

You can either specify the project name directly as an argument or use the interactive prompt to select a project.
Use the '--output-format' flag for raw JSON output.

```sh
harbor webhook list [PROJECT_NAME] [flags]
```

### Examples

```sh
  # List webhooks for a specific project
  harbor-cli webhook list my-project

  # List webhooks interactively by selecting the project
  harbor-cli webhook list

  # Output in JSON format
  harbor-cli webhook list my-project --output-format=json
```

### Options

```sh
  -h, --help   help for list
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor webhook](harbor-webhook.md)	 - Manage webhook policies in Harbor

