# harbor-cli
LFX mentorship program project original author https://github.com/akshatdalton 

https://github.com/cncf/mentoring/tree/main/programs/lfx-mentorship/2023/01-Mar-May#an-official-golang-api-client-and-cli-for-harbor
https://mentorship.lfx.linuxfoundation.org/project/7e8cb88a-5b37-471c-8db8-e11907b5a661

#Under Contstruction! 
# How to Install:

- Run `go mod download` to install the dependencies
- Run `go build -o harbor` to generate the build file

# How to Use:

If you have completed the above steps, now you are all set to use this project.
<br>
<br>
`./harbor --help` or do `./harbor help`:
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
