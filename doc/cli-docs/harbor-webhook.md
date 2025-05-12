---
title: harbor webhook
weight: 0
---
## harbor webhook

### Description

##### Manage webhook policies in Harbor

### Synopsis

Use this command to manage webhook policies in your Harbor projects.

Webhooks enable external systems to be notified of events in Harbor (e.g., pushing an image, deleting an artifact). 
This command supports listing, creating, editing, and deleting webhook configurations.

### Examples

```sh
  # List webhook policies in a project
  harbor-cli webhook list my-project

  # Create a new webhook policy
  harbor-cli webhook create --project my-project --name my-webhook

  # Edit an existing webhook policy
  harbor-cli webhook edit --project my-project --webhook-id 5

  # Delete a webhook policy
  harbor-cli webhook delete --project my-project --webhook-id 5
```

### Options

```sh
  -h, --help   help for webhook
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI
* [harbor webhook create](harbor-webhook-create.md)	 - Create a new webhook for a Harbor project
* [harbor webhook delete](harbor-webhook-delete.md)	 - Delete a webhook from a Harbor project
* [harbor webhook edit](harbor-webhook-edit.md)	 - Edit an existing webhook for a Harbor project
* [harbor webhook list](harbor-webhook-list.md)	 - List all webhook policies for a Harbor project

