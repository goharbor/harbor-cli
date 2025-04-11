---
title: harbor webhook delete
weight: 85
---
## harbor webhook delete

### Description

##### Delete a webhook from a Harbor project

### Synopsis

This command deletes a webhook from the specified Harbor project.
You can either specify the project name and webhook ID directly using flags,
or interactively select a project and webhook if not provided.

```sh
harbor webhook delete [flags]
```

### Examples

```sh
  # Delete a webhook by specifying the project and webhook ID
  harbor-cli webhook delete --project my-project --webhook 5

  # Delete a webhook by selecting the project and webhook interactively
  harbor-cli webhook delete
```

### Options

```sh
  -h, --help             help for delete
      --project string   Project Name
      --webhook string   Webhook ID
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor webhook](harbor-webhook.md)	 - Manage webhook policies in Harbor

