---
title: harbor project config update
weight: 10
---
## harbor project config update

### Description

##### Interactively or via flags update project configuration in Harbor

### Synopsis

Update the configuration settings of a Harbor project either interactively or directly using command-line flags.

You can specify the project by its name or ID as an argument. If not provided, you will be prompted to select a project interactively.

Examples:

  # Update project 'myproject' visibility to public
  harbor-cli project config update myproject --public true

  # Update multiple settings in one command
  harbor-cli project config update myproject --public false --prevent-vul true --severity high

  # Run interactively without flags
  harbor-cli project config update

Supported flag values:

  - Boolean flags (public, auto-scan, prevent-vul, reuse-sys-cve-allowlist, enable-content-trust, enable-content-trust-cosign): "true" or "false"
  - Severity: one of "low", "medium", "high", "critical"


```sh
harbor project config update [project_name] [flags]
```

### Options

```sh
      --auto-scan string                     Enable or disable auto scan (true/false)
      --enable-content-trust string          Enable or disable content trust (true/false)
      --enable-content-trust-cosign string   Enable or disable content trust cosign (true/false)
  -h, --help                                 help for update
      --prevent-vul string                   Enable or disable vulnerability prevention (true/false)
      --public string                        Set project visibility (true/false)
      --reuse-sys-cve string                 Enable or disable reuse of system CVE allowlist (true/false)
      --severity string                      Set severity level
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
      --id                     Use project ID instead of name
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project config](harbor-project-config.md)	 - Manage project configuration

