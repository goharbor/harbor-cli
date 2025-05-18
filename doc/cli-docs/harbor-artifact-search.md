---
title: harbor artifact search
weight: 5
---
## harbor artifact search

### Description

##### Search container artifacts (images, charts, etc.) in a Harbor repository

### Synopsis

List artifacts (e.g., container images, charts) within a given Harbor repository. 
Search is based on matching tags and artifact types (e.g., container, images, charts)

Examples:
  harbor-cli artifact search project/repo:tag               
  harbor-cli artifact search project/repo:tag --type IMAGE


```sh
harbor artifact search [flags]
```

### Options

```sh
  -h, --help            help for search
  -p, --page int        Page number (default 1)
  -n, --page-size int   Size of per page (default 10)
  -q, --query string    Query string to query resources
  -s, --sort string     Sort the resource list in ascending or descending order
  -t, --type string     Filter artifacts by type (e.g., IMAGE, CHART)
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor artifact](harbor-artifact.md)	 - Manage artifacts

