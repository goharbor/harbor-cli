---
title: harbor scan all update schedule
weight: 50
---
## harbor scan-all update-schedule

### Description

##### update-schedule [schedule-type: none|hourly|daily|weekly|custom]

### Synopsis

Configure or update the automatic vulnerability scan schedule for all artifacts.

This command allows you to set when Harbor automatically scans all artifacts for vulnerabilities. You can choose from predefined schedules or create a custom schedule using cron expressions.

Available schedule types:
  - none:    Disable automatic scanning
  - hourly:  Run scan every hour
  - daily:   Run scan once per day
  - weekly:  Run scan once per week
  - custom:  Define a custom schedule using a cron expression

For custom schedules, Harbor requires a 6-field cron expression in the format:
  seconds minutes hours day-of-month month day-of-week

Examples:
  # Disable scheduled scanning
  harbor-cli scan-all update-schedule none

  # Set daily automatic scanning
  harbor-cli scan-all update-schedule daily

  # Set weekly automatic scanning
  harbor-cli scan-all update-schedule weekly

  # Set a custom schedule (every day at 2:30 AM)
  harbor-cli scan-all update-schedule custom --cron "0 30 2 * * *"

  # Use interactive mode to configure a custom schedule
  harbor-cli scan-all update-schedule custom

Note: For custom schedules, if you provide a 5-field cron expression, the CLI will automatically add a leading "0" for the seconds field to create the required 6-field format.

```sh
harbor scan-all update-schedule [flags]
```

### Options

```sh
      --cron string   Cron expression for custom schedule (include the expression in double quotes)
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

