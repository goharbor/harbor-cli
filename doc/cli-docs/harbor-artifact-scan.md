---
title: harbor artifact scan
weight: 60
---
## harbor artifact scan

### Description

##### Scan an artifact

### Synopsis

Scan an artifact in Harbor Repository

### Examples

```sh
harbor artifact scan start <project>/<repository>/<reference>
```

### Options

```sh
  -h, --help   help for scan
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor artifact](harbor-artifact.md)	 - Manage artifacts
* [harbor artifact scan start](harbor-artifact-scan-start.md)	 - Start a scan of an artifact
* [harbor artifact scan stop](harbor-artifact-scan-stop.md)	 - Stop a scan of an artifact

