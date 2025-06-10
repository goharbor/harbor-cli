---
title: harbor project robot refresh
weight: 25
---
## harbor project robot refresh

### Description

##### refresh robot secret by id

### Synopsis

Refresh the secret for an existing robot account in Harbor.

This command generates a new secret for a robot account, effectively revoking 
the old secret and requiring updates to any systems using the robot's credentials.

The command supports multiple ways to identify the robot account:
- By providing the robot ID directly as an argument
- Without any arguments, which will prompt for both project and robot selection

You can specify the new secret in several ways:
- Let Harbor generate a random secret (default)
- Provide a custom secret with the --secret flag
- Pipe a secret via stdin using the --secret-stdin flag

After refreshing, the new secret will be:
- Displayed on screen
- Copied to clipboard for immediate use
- Usable immediately for authentication

Important considerations:
- The old secret will be invalidated immediately
- Any systems using the old credentials will need to be updated
- There is no way to recover the old secret after refreshing

Examples:
  # Refresh robot secret by ID (generates a random secret)
  harbor-cli project robot refresh 123

  # Refresh with a custom secret
  harbor-cli project robot refresh 123 --secret "MyCustomSecret123"

  # Provide secret via stdin (useful for scripting)
  echo "MySecretFromScript123" | harbor-cli project robot refresh 123 --secret-stdin

  # Interactive refresh (will prompt for project and robot selection)
  harbor-cli project robot refresh

```sh
harbor project robot refresh [robotID] [flags]
```

### Options

```sh
  -h, --help            help for refresh
      --secret string   secret
      --secret-stdin    Take the robot secret from stdin
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project robot](harbor-project-robot.md)	 - Manage robot accounts

