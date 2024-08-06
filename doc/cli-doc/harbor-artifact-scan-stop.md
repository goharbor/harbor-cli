---
title: harbor artifact scan stop
weight: 40
---
## harbor artifact scan stop

#### Stop a scan of an artifact

### Synopsis

Stop a scan of an artifact in Harbor Repository

```sh
harbor artifact scan stop [flags]
```

### Examples

```sh
harbor artifact scan stop <project>/<repository>/<reference>
```

### Options

```sh
  -h, --help   help for stop
```

### Options inherited from parent commands

```sh
      --config string          config file (default is $HOME/.harbor/config.yaml) (default "/home/user/.harbor/config.yaml")
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

