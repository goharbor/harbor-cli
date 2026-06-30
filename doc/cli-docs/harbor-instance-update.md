---
title: harbor instance update
weight: 20
---
## harbor instance update

### Description

##### Update a preheat provider instance in Harbor

### Synopsis

Update a preheat provider instance in Harbor by name or ID. If no update
flags are provided, the command opens an interactive update form.

```sh
harbor instance update [NAME|ID] [flags]
```

### Examples

```sh
  harbor-cli instance update my-instance --description "Updated preheat instance"
  harbor-cli instance update 1 --id --enable=false
  harbor-cli instance update my-instance --authmode BASIC --auth-username admin --auth-password Harbor12345
```

### Options

```sh
      --auth-password string   Password for BASIC authentication
      --auth-token string      Token for OAUTH authentication
      --auth-username string   Username for BASIC authentication
  -a, --authmode string        Authentication mode (NONE, BASIC, OAUTH)
  -d, --description string     Description of the instance
      --enable                 Whether the instance is enabled or not
  -h, --help                   help for update
      --id                     Get instance by id
  -i, --insecure               Whether or not the certificate will be verified when Harbor tries to access the server
  -n, --name string            New name for the instance
  -u, --url string             Endpoint URL for the instance
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor instance](harbor-instance.md)	 - Manage preheat provider instances in Harbor

