---
title: harbor artifact list
weight: 20
---
## harbor artifact list

### Description

##### List container artifacts (images, charts, etc.) in a Harbor repository with metadata

### Synopsis

List all artifacts (e.g., container images, charts) within a given Harbor repository. 
Supports optional project/repository input in the form <project>/<repository>. 
Displays key artifact metadata including tags, digest, type, size, vulnerability count, and push time.

Examples:
  harbor-cli artifact list                # Interactive prompt for project and repository
  harbor-cli artifact list library/nginx # Directly list artifacts in the nginx repo under 'library' project

Supports pagination, search queries, and sorting using flags.

```sh
harbor artifact list [flags]
```

### Options

```sh
  -h, --help            help for list
  -p, --page int        Page number (default 1)
  -n, --page-size int   Size of per page (default 10)
  -q, --query string    Query string to query resources
  -s, --sort string     Sort the resource list in ascending or descending order
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor artifact](harbor-artifact.md)	 - Manage artifacts

