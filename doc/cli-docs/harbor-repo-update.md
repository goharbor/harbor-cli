---
title: harbor repo update
weight: 55
---
## harbor repo update

### Description

##### Update a repository

### Synopsis

Update the description of a repository.

This command updates the description associated with a repository
within a Harbor project.

Examples:
  # Update repository description using project/repository format
  	harbor repository update library/nginx --description "Official nginx image"

```sh
harbor repo update [flags]
```

### Options

```sh
  -d, --description string   Repository description
  -h, --help                 help for update
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor repo](harbor-repo.md)	 - Manage repositories

