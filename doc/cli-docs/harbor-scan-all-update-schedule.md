---
title: harbor scan all update schedule
weight: 50
---
## harbor scan-all update-schedule

### Description

##### update-schedule [schedule-type: none|hourly|daily|weekly|custom]

```sh
harbor scan-all update-schedule [flags]
```

### Options

```sh
      --cron string   Cron expression (include the expression in double quotes)
  -h, --help          help for update-schedule
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor scan-all](harbor-scan-all.md)	 - Scan all artifacts

