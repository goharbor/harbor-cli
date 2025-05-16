---
title: harbor webhook delete
weight: 70
---
## harbor webhook delete

### Description

##### Delete a webhook from a Harbor project

### Synopsis

Delete a webhook from a specified Harbor project.
You can either specify the project name and webhook ID using flags,
pass the webhook name as an argument,
or interactively select a project and webhook if not provided.

```sh
harbor webhook delete [webhook-name] [flags]
```

### Examples

```sh
  # Delete by project and webhook ID
  harbor-cli webhook delete --project my-project --webhook 5

  # Delete by project and webhook name
  harbor-cli webhook delete my-webhook --project my-project

  # Fully interactive deletion
  harbor-cli webhook delete
```

### Options

```sh
  -h, --help             help for delete
      --project string   Project name (required when providing webhook ID or name)
      --webhook-id int   Webhook ID (alternative to providing webhook name as argument) (default -1)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor webhook](harbor-webhook.md)	 - Manage webhook policies in Harbor

