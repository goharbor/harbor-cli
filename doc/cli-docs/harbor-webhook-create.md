---
title: harbor webhook create
weight: 60
---
## harbor webhook create

### Description

##### Create a new webhook for a Harbor project

### Synopsis

This command creates a new webhook policy for a specified Harbor project.

You can either provide all required flags (project name, notify type, endpoint, etc.) directly to create the webhook non-interactively,
or leave them out and be guided through an interactive prompt to input each field. The webhook name is required as an argument.

```sh
harbor webhook create [name] [flags]
```

### Examples

```sh
  # Create a webhook using flags
  harbor-cli webhook create my-webhook \
    --project my-project \
    --notify-type http \
    --event-type PUSH_ARTIFACT,DELETE_ARTIFACT \
    --endpoint-url https://example.com/webhook \
    --description "Webhook for artifact events" \
    --payload-format Default \
    --auth-header "Bearer mytoken"

  # Create a webhook using the interactive prompt
  harbor-cli webhook create my-webhook
```

### Options

```sh
      --auth-header string          Authentication Header
      --description string          Webhook Description
      --endpoint-url string         Webhook Endpoint URL
      --event-type strings          Event Types (comma separated)
  -h, --help                        help for create
      --notify-type string          Notify Type (http, slack)
      --payload-format string       Payload Format (Default, CloudEvents)
      --project string              Project Name
      --verify-remote-certificate   Verify Remote Certificate (default true)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor webhook](harbor-webhook.md)	 - Manage webhook policies in Harbor

