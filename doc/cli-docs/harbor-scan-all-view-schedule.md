---
title: harbor scan all view schedule
weight: 35
---
## harbor scan-all view-schedule

### Description

##### View the scan all schedule

### Synopsis

Display the current vulnerability scan schedule configuration.

This command retrieves and shows the current automatic scanning schedule settings for your Harbor instance, including:

- Schedule Type: The type of schedule (None, Hourly, Daily, Weekly, or Custom)
- Cron Expression: For custom schedules, shows the configured cron pattern
- Next Scheduled Time: When the next automatic scan is scheduled to run

This information helps you understand when Harbor will automatically scan your artifacts
for vulnerabilities.

Examples:
  # View the current scan schedule
  harbor-cli scan-all view-schedule

  # View the schedule in JSON format
  harbor-cli scan-all view-schedule --output-format json

You can use this command to verify changes after updating the schedule with the 'update-schedule' command.

```sh
harbor scan-all view-schedule [flags]
```

### Options

```sh
  -h, --help   help for view-schedule
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor scan-all](harbor-scan-all.md)	 - Scan all artifacts

