# Harbor CLI Replication Policy Configuration

This guide explains how to use YAML or JSON configuration files to create Harbor replication policies with the Harbor CLI.

## Overview

Instead of using interactive prompts, you can define replication policies in configuration files and apply them using the `--policy-config-file` or `-f` flag.

## Usage

```bash
# Using YAML file
./harbor-cli replication policies create -f policy.yaml

# Using JSON file
./harbor-cli replication policies create -f policy.json
```

## Configuration File Structure

### Basic Policy Configuration

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Unique name for the replication policy |
| `description` | string | ❌ | Optional description of the policy |
| `replication_mode` | string | ✅ | `"push"` or `"pull"` |
| `target_registry` | string | ✅ | Name of the target registry |
| `trigger_mode` | string | ✅ | `"manual"`, `"scheduled"`, or `"event_based"` |
| `cron_string` | string | ❌ | Required for scheduled triggers (6-field cron format) |
| `bandwidth_limit` | string | ❌ | Bandwidth limit in KB/s (`"-1"` for unlimited) |
| `override` | boolean | ❌ | Replace existing artifacts (default: `false`) |
| `replicate_deletion` | boolean | ❌ | Sync deletion operations (default: `false`) |
| `copy_by_chunk` | boolean | ❌ | Transfer in chunks (default: `false`) |
| `enabled` | boolean | ❌ | Enable policy after creation (default: `false`) |

### Replication Filters

The `replication_filter` array contains filter objects that determine which artifacts to replicate.

#### Filter Types

| Type | Description | Decoration Support | Valid Values |
|------|-------------|-------------------|--------------|
| `resource` | Artifact type | ❌ | `"image"` |
| `name` | Repository name pattern | ❌ | Any string with wildcards (`*`) |
| `tag` | Tag pattern | ✅ | Any string with wildcards (`*`) |
| `label` | Label key=value pattern | ✅ | `"key=value"` format |

#### Decoration Values
- `"matches"` - Include artifacts matching the pattern
- `"excludes"` - Exclude artifacts matching the pattern

**Note:** Only `tag` and `label` filters support decoration. `resource` and `name` filters always match.

## Examples

### 1. Manual Pull Policy (YAML)

```yaml
# manual-pull.yaml
name: "manual-pull-from-dockerhub"
description: "Manually pull specific nginx images from Docker Hub"
replication_mode: "pull"
target_registry: "dockerhub-proxy"
trigger_mode: "manual"
bandwidth_limit: "-1"
override: false
replicate_deletion: false
enabled: true
copy_by_chunk: false

replication_filter:
  - type: "resource"
    value: "image"
  
  - type: "name"
    value: "nginx"
  
  - type: "tag"
    decoration: "matches"
    value: "1.*"
  
  - type: "label"
    decoration: "matches"
    value: "maintainer=NGINX Docker Maintainers"
```

### 2. Scheduled Push Policy (JSON)

```json
{
  "name": "nightly-backup-push",
  "description": "Push all production images to backup registry every night",
  "replication_mode": "push",
  "target_registry": "backup-harbor",
  "trigger_mode": "scheduled",
  "cron_string": "0 0 2 * * *",
  "bandwidth_limit": "2048",
  "override": true,
  "replicate_deletion": true,
  "enabled": true,
  "copy_by_chunk": true,
  "replication_filter": [
    {
      "type": "resource",
      "value": "image"
    },
    {
      "type": "name",
      "value": "production/*"
    },
    {
      "type": "tag",
      "decoration": "excludes",
      "value": "*-snapshot"
    },
    {
      "type": "tag",
      "decoration": "excludes",
      "value": "*-dev"
    },
    {
      "type": "label",
      "decoration": "matches",
      "value": "release=stable"
    }
  ]
}
```

### 3. Event-Based Push Policy (YAML)

```yaml
# event-based.yaml
name: "realtime-sync"
description: "Automatically sync artifacts to partner registry on push"
replication_mode: "push"
target_registry: "partner-registry"
trigger_mode: "event_based"
bandwidth_limit: "512"
override: true
replicate_deletion: true
enabled: true
copy_by_chunk: false

replication_filter:
  - type: "name"
    value: "apps/*"
  
  - type: "tag"
    decoration: "matches"
    value: "[0-9]*.[0-9]*.[0-9]*"
  
  - type: "label"
    decoration: "excludes"
    value: "stage=development"
  
  - type: "label"
    decoration: "matches"
    value: "public=true"
```

## Trigger Modes

### Manual
- Replication runs only when manually triggered
- No additional configuration required

### Scheduled
- Replication runs on a cron schedule
- Requires `cron_string` field with 6-field format: `"seconds minutes hours day-month month day-week"`
- Example: `"0 0 2 * * *"` (daily at 2:00 AM)

### Event-Based
- Replication runs automatically when artifacts are pushed/updated
- Only available for `"push"` mode
- Supports `replicate_deletion` for syncing deletions

## Common Patterns

### Filter by Repository Namespace
```yaml
replication_filter:
  - type: "name"
    value: "library/*"  # All repositories in 'library' namespace
```

### Version-Specific Filtering
```yaml
replication_filter:
  - type: "tag"
    decoration: "matches"
    value: "v*"  # Only tags starting with 'v'
  
  - type: "tag"
    decoration: "excludes"
    value: "*-rc*"  # Exclude release candidates
```

### Label-Based Filtering
```yaml
replication_filter:
  - type: "label"
    decoration: "matches"
    value: "environment=production"
  
  - type: "label"
    decoration: "excludes"
    value: "experimental=true"
```

## Validation Rules

- Policy `name` must be unique and non-empty
- `replication_mode` must be `"push"` or `"pull"`
- `trigger_mode` must be `"manual"`, `"scheduled"`, or `"event_based"`
- `cron_string` is required for scheduled triggers
- `event_based` triggers are only available for push mode
- Resource filter value must be `"image"` or  `"artifact"` (if specified)
- Tag and label filters support decoration, others don't

## File Formats

Both YAML and JSON formats are supported:
- `.yaml` or `.yml` files are parsed as YAML
- `.json` files are parsed as JSON
- File extension is required

## Error Handling

Common validation errors:
- `name is required` - Missing policy name
- `replication_mode must be 'push' or 'pull'` - Invalid replication mode
- `decoration is only supported for 'tag' and 'label' filters` - Invalid decoration usage
- `resource value must be 'image'` - Invalid resource filter value
- `cron string cannot be empty for scheduled trigger` - Missing cron string

## Best Practices

1. **Use descriptive names** - Policy names should clearly indicate their purpose
2. **Start with manual triggers** - Test policies manually before automating
3. **Set bandwidth limits** - Prevent replication from overwhelming network
4. **Use specific filters** - Avoid replicating unnecessary artifacts
5. **Enable deletion sync carefully** - Only use when you want true synchronization
6. **Test with small datasets** - Validate filters with limited scope first

## Registry Configuration

Before using replication policies, ensure your target registries are configured in Harbor:

```bash
# List available registries
./harbor-cli registry list

# Create a new registry if needed
./harbor-cli registry create
```

The `target_registry` field in your configuration must match an existing registry name in Harbor.