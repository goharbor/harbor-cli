---
title: harbor tag immutable list
weight: 40
---
## harbor tag immutable list

### Description

##### Display all immutable tag rules for a project

### Synopsis

Retrieve and display a list of immutable tag rules configured for a specified project in Harbor. 
Immutable tag rules prevent specific tags from being deleted or overwritten, ensuring better security and compliance.
You can specify the project name as an argument or, if omitted, you will be prompted to select one interactively.

```sh
harbor tag immutable list [PROJECT_NAME] [flags]
```

### Examples

```sh
  
  # List immutable tag rules for a specific project  
  harbor tag immutable list my-project  

  # List immutable tag rules interactively (if no project name is provided)  
  harbor tag immutable list  
  
```

### Options

```sh
  -h, --help   help for list
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor tag immutable](harbor-tag-immutable.md)	 - Manage Immutability rules in the project

