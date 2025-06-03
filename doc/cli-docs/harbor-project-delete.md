---
title: harbor project delete
weight: 100
---
## harbor project delete

### Description

##### Delete project by name or ID

### Synopsis

Delete project by name or ID. Multiple projects can be deleted by providing their names as arguments. If no arguments are provided, it will prompt for the project name. Use --project-id to specify the project ID for single project directly. The --force flag will delete all repositories and artifacts within the project.

```sh
harbor project delete [flags]
```

### Examples

```sh
harbor project delete [projectname1] [projectname2] or harbor project delete --project-id [projectid]
```

### Options

```sh
      --force               Forcefully delete all repositories, artifacts, and policies in the project. Use with extreme caution—this action is irreversible.
  -h, --help                help for delete
      --project-id string   Specify project ID instead of project name
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor project](harbor-project.md)	 - Manage projects and assign resources to them

