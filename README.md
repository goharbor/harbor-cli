# harbor-cli
The cli for harbor registry. This project is under construction.

### LFX mentorship program

LFX mentorship program project original author [Akshat](https://github.com/akshatdalton) 

LFX mentorship term 2 project proposal: [Harbor CLI](https://github.com/cncf/mentoring/tree/main/programs/lfx-mentorship/2023/01-Mar-May#an-official-golang-api-client-and-cli-for-harbor)
<br>
LFX mentorship website: [LFX mentorship Harbor CLI](https://mentorship.lfx.linuxfoundation.org/project/7e8cb88a-5b37-471c-8db8-e11907b5a661)

### Installation
- Run `go mod download` to install the dependencies
- Run `go build -o harbor` to generate the build file

### Usage:

`./harbor --help` or `./harbor help` to get the list of commands:
```shell
Official Harbor CLI

Usage:
  harbor [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      create project, registry, etc.
  delete      delete project, registry, etc.
  get         get project, registry, etc.
  help        Help about any command
  list        list project, registry, etc.
  login       Log in to Harbor registry
  update      update registry, etc.

Flags:
  -h, --help   help for harbor

Use "harbor [command] --help" for more information about a command.
```
