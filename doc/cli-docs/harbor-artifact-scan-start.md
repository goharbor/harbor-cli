---
title: harbor artifact scan start
weight: 80
---
## harbor artifact scan start

### Description

##### Start a scan of an artifact

### Synopsis

Start a scan of an artifact in Harbor Repository

```sh
harbor artifact scan start [flags]
```

### Examples

```sh
harbor artifact scan start <project>/<repository>/<reference>
```

### Options

```sh
  -h, --help   help for start
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor artifact scan](harbor-artifact-scan.md)	 - Scan an artifact

