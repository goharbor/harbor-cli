---
title: harbor webhook edit
weight: 25
---
## harbor webhook edit

### Description

##### webhook edit

```sh
harbor webhook edit [flags]
```

### Options

```sh
      --auth-header string          Authentication Header
      --description string          Webhook Description
      --enabled                     Webhook Enabled (default true)
      --endpoint-url string         Webhook Endpoint URL
      --event-type stringArray      Event Types (comma separated)
  -h, --help                        help for edit
      --name string                 Webhook Name
      --notify-type string          Notify Type (http, slack)
      --payload-format string       Payload Format (Default, CloudEvents)
      --project string              Project Name
      --verify-remote-certificate   Verify Remote Certificate (default true)
      --webhook-id string           Webhook ID
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor webhook](harbor-webhook.md)	 - Manage webhooks

