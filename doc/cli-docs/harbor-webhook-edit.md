---
title: harbor webhook edit
weight: 55
---
## harbor webhook edit

### Description

##### Edit an existing webhook for a Harbor project

### Synopsis

This command allows you to update an existing webhook policy in a Harbor project.

You can either pass all the necessary flags (webhook ID, project name, etc.) to perform a non-interactive update,
or leave them out and use the interactive prompt to select and update a webhook.

```sh
harbor webhook edit [WEBHOOK_NAME] [flags]
```

### Examples

```sh
  # Edit a webhook by providing all fields directly
  harbor-cli webhook edit my-webhook \
    --project my-project \
    --notify-type http \
    --event-type PUSH_ARTIFACT \
    --endpoint-url https://new-url.com \
    --description "Updated webhook for artifact push" \
    --payload-format Default \
    --auth-header "Bearer newtoken" \
    --enabled=true

  # Edit a webhook using the interactive prompt
  harbor-cli webhook edit
```

### Options

```sh
      --auth-header string          Authentication Header
      --description string          Webhook Description
      --enabled                     Webhook Enabled (default true)
      --endpoint-url string         Webhook Endpoint URL
      --event-type strings          Event Types (comma separated)
  -h, --help                        help for edit
      --notify-type string          Notify Type (http, slack)
      --payload-format string       Payload Format (Default, CloudEvents)
      --project string              Project Name
      --verify-remote-certificate   Verify Remote Certificate (default true)
      --webhook-id int              Webhook ID (default -1)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor webhook](harbor-webhook.md)	 - Manage webhook policies in Harbor

