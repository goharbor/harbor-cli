---
title: harbor cve allowlist
weight: 70
---
## harbor cve-allowlist

### Description

##### Manage system CVE allowlist

### Synopsis

Managing CVE lists that are intentionally excluded from vulnerability scanning

### Examples

```sh
harbor cve-allowlist list
```

### Options

```sh
  -h, --help   help for cve-allowlist
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor](harbor.md)	 - Official Harbor CLI
* [harbor cve-allowlist add](harbor-cve-allowlist-add.md)	 - Add cve allowlist
* [harbor cve-allowlist list](harbor-cve-allowlist-list.md)	 - List system level allowlist of cve

