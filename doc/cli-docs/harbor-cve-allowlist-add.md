---
title: harbor cve allowlist add
weight: 60
---
## harbor cve-allowlist add

### Description

##### Add cve allowlist

### Synopsis

Create allowlist of CVEs to ignore during vulnerability scanning

```sh
harbor cve-allowlist add [flags]
```

### Options

```sh
  -n, --cveid string        Comma-separated list of CVE IDs to be added to the allowlist
  -d, --expiredate string   Specifies the expiration date for the CVE entries in the format 'YYYY-MM-DD'
  -h, --help                help for add
  -i, --isexpire            Indicates whether the CVE entries should have an expiration date. Set to true to specify an expiration date
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor cve-allowlist](harbor-cve-allowlist.md)	 - Manage system CVE allowlist

