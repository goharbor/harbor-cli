---
title: harbor artifact scan start
weight: 35
---
## harbor artifact scan start

#### Start a scan of an artifact

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
      --config string          config file (default is $HOME/.harbor/config.yaml) (default "/home/user/.harbor/config.yaml")
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

