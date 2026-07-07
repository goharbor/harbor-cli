---
title: harbor project update
weight: 82
---
## harbor project update

### Description

##### Update a project

### Synopsis

Update project settings such as visibility, storage quota, and metadata.

Examples:
  harbor project update myproject --public true
  harbor project update myproject --storage-limit -1 --prevent-vul true

```sh
harbor project update [project_name] [flags]
```

### Options

```sh
      --auto-scan string       Enable or disable auto scan (true/false)
  -h, --help                   help for update
      --id                     Use project ID instead of name
      --prevent-vul string     Enable or disable vulnerability prevention (true/false)
      --public string          Set project visibility (true/false)
      --registry-id string     ID of referenced registry for proxy cache projects
      --reuse-sys-cve string   Enable or disable reuse of system CVE allowlist (true/false)
      --severity string        Set severity level (none, low, medium, high, critical)
      --storage-limit string   Storage quota of the project in bytes (-1 for unlimited)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project](harbor-project.md)	 - Manage projects and assign resources to them

