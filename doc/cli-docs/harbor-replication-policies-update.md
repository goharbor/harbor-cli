---
title: harbor replication policies update
weight: 5
---
## harbor replication policies update

### Description

##### Update an existing replication policy

When called without any update flags, the command opens the interactive TUI wizard
(existing behavior preserved). When any update flag is provided, the command runs
non-interactively — loading the existing policy as the baseline and applying only the
explicitly provided flags (partial update).

```sh
harbor replication policies update [policy-id] [flags]
```

### Examples

```sh
  # Interactive (existing behavior — opens TUI wizard)
  harbor replication policies update 1

  # Change only the name
  harbor replication policies update 1 --name production-sync

  # Enable the policy
  harbor replication policies update 1 --enabled=true

  # Update multiple fields at once
  harbor replication policies update 1 \
    --name production-sync \
    --description "Production replication" \
    --enabled=true \
    --override=true

  # Update resource/name/tag filters
  harbor replication policies update 1 \
    --name-filter "library/" \
    --tag-filter matches \
    --tag-pattern "v*"

  # Switch to a scheduled trigger
  harbor replication policies update 1 \
    --trigger-type scheduled \
    --cron "0 0 */6 * * *"
```

### Options

```sh
      --copy-by-chunk            Transfer artifacts in chunks for better reliability
      --cron string              Cron schedule (6-field format, required when --trigger-type=scheduled, e.g. "0 0 */6 * * *")
      --description string       New description for the replication policy
      --enabled                  Enable the replication policy
  -h, --help                     help for update
      --label-filter string      Label filter type: matches or excludes
      --label-pattern string     Label filter pattern (e.g. env=prod or env=prod,ver=1.0)
      --name string              New name for the replication policy
      --name-filter string       Repository name filter pattern (supports wildcards, e.g. library/*)
      --override                 Override artifacts on destination if they already exist
      --replicate-deletion       Replicate deletion operations to the destination
      --resource-filter string   Resource type filter: image, artifact, or empty for all
      --speed string             Maximum replication speed in KB/s (-1 for unlimited)
      --tag-filter string        Tag filter type: matches or excludes
      --tag-pattern string       Tag filter pattern (e.g. v*, latest, *-prod)
      --trigger-type string      Trigger type: manual, scheduled, or event_based
```

### Options inherited from parent commands

```sh
  -c, --config string          config file (default is $HOME/.config/harbor-cli/config.yaml)
  -o, --output-format string   Output format. One of: json|yaml|csv
  -v, --verbose                verbose output
```

### SEE ALSO

* [harbor replication policies](harbor-replication-policies.md)	 - Manage replication policies

