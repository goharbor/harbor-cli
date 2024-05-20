# harbor-cli
[![Go Report Card](https://goreportcard.com/badge/github.com/goharbor/harbor-cli)](https://goreportcard.com/report/github.com/goharbor/harbor-cli)
<br>

The cli for harbor registry. This project is under construction.

|Community Meeting|
|------------------|
|The Harbor Project holds bi-weekly community calls in two different timezones. To join the community calls or to watch previous meeting notes and recordings, please visit the [meeting schedule](https://github.com/goharbor/community/blob/master/MEETING_SCHEDULE.md).|

### LFX mentorship program

LFX mentorship program project original author [Akshat](https://github.com/akshatdalton) 

LFX mentorship term 2 project proposal: [Harbor CLI](https://github.com/cncf/mentoring/tree/main/programs/lfx-mentorship/2023/01-Mar-May#an-official-golang-api-client-and-cli-for-harbor)
<br>
LFX mentorship website: [LFX mentorship Harbor CLI project](https://mentorship.lfx.linuxfoundation.org/project/7e8cb88a-5b37-471c-8db8-e11907b5a661)

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

### Architecture

For learning the architecture design of Harbor, check the document [Architecture Overview of Harbor](https://github.com/goharbor/harbor/wiki/Architecture-Overview-of-Harbor).

### API

* Harbor RESTful API: The APIs for most administrative operations of Harbor and can be used to perform integrations with Harbor programmatically.
    * Part 1: [New or changed APIs](https://editor.swagger.io/?url=https://raw.githubusercontent.com/goharbor/harbor/main/api/v2.0/swagger.yaml)

### OCI Distribution Conformance Tests

Check the OCI distribution conformance tests [report](https://storage.googleapis.com/harbor-conformance-test/report.html) of Harbor.

### Community

* **Twitter:** [@project_harbor](https://twitter.com/project_harbor)
* **User Group:** Join Harbor user email group: [harbor-users@lists.cncf.io][users-dl] to get update of Harbor's news, features, releases, or to provide suggestion and feedback.
* **Developer Group:** Join Harbor developer group: [harbor-dev@lists.cncf.io][dev-dl] for discussion on Harbor development and contribution.
* **Slack:** Join Harbor's community for discussion and ask questions: [Cloud Native Computing Foundation][cncf-slack], channel: [#harbor][users-slack] and [#harbor-dev][dev-slack]
* [Bi-weekly public community meetings][community-meetings]
  * Catch up with [past meetings on YouTube][past-meetings]

### Contributing
* Fork and clone the repo
* Create a new branch for the changes
* Make your changes on that branch
* Raise a PR against the main branch of the main repo
* The PR will be reviewed by the maintainers
* If necessary, the maintainers will ask for changes
* Once the PR is approved, it will be merged

#### Important
* Keep sync with the main branch. If your branch is behind the main branch, please rebase your branch to the main branch. Please do not use git pull as it leaves merge commits, which makes the commit history messy.
<br>

* As Harbor has integrated the [DCO (Developer Certificate of Origin)](https://probot.github.io/apps/dco/) check tool, contributors are required to sign off that they adhere to those requirements by adding a `Signed-off-by` line to the commit messages. Git has even provided a `-s` command line option to append that automatically to your commit messages, please use it when you commit your changes.

```bash
$ git commit -s -m 'This is my commit message'
```

### Demos

* **[Live Demo](https://demo.goharbor.io)** - A demo environment with the latest Harbor stable build installed. For additional information please refer to [this page](https://goharbor.io/docs/latest/install-config/demo-server/).
* **[Video Demos](https://github.com/goharbor/harbor/wiki/Video-demos-for-Harbor)** - Demos for Harbor features and continuously updated.

### Security

#### Security Audit

A third party security audit was performed by Cure53 in October 2019. You can see the full report [here](https://goharbor.io/docs/2.0.0/security/Harbor_Security_Audit_Oct2019.pdf).

#### Reporting security vulnerabilities

If you've found a security related issue, a vulnerability, or a potential vulnerability in Harbor please let the [Harbor Security Team](mailto:cncf-harbor-security@lists.cncf.io) know with the details of the vulnerability. We'll send a confirmation
email to acknowledge your report, and we'll send an additional email when we've identified the issue
positively or negatively.

For further details please see our complete [security release process][harbor-security].






[community-meetings]: https://github.com/goharbor/community/blob/main/MEETING_SCHEDULE.md
[past-meetings]: https://www.youtube.com/playlist?list=PLgInP-D86bCwTC0DYAa1pgupsQIAWPomv
[users-slack]: https://cloud-native.slack.com/archives/CC1E09J6S
[dev-slack]: https://cloud-native.slack.com/archives/CC1E0J0MC
[cncf-slack]: https://slack.cncf.io
[users-dl]: https://lists.cncf.io/g/harbor-users
[dev-dl]: https://lists.cncf.io/g/harbor-dev
[twitter]: http://twitter.com/project_harbor
[harbor-security]: https://github.com/goharbor/harbor/blob/main/SECURITY.md