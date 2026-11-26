# Declarative Harbor configuration

The declarative document describes desired Harbor resources rather than Harbor
API responses. Generated IDs, timestamps, resource usage, health status, and
credentials are intentionally excluded.

The initial `goharbor.io/v1alpha1` contract supports:

- Harbor system configuration
- external registry endpoints
- projects, project metadata, quotas, and webhooks
- replication policies

Resource names are their portable identities. References such as a project's
proxy-cache registry and a replication policy's external registry therefore use
names instead of Harbor-generated IDs.

Export current state:

```sh
harbor export -f harbor.yaml
```

Review and apply it:

```sh
harbor apply -f harbor.yaml --dry-run
harbor apply -f harbor.yaml
```

Configurations can also be organized as a directory:

```text
production/
├── 00-base.yaml
├── 10-production.yaml
└── projects/
    └── 20-applications.yaml
```

```sh
harbor apply -f production --dry-run
harbor apply -f production
```

The CLI discovers YAML and JSON files recursively and merges them in
lexicographic relative-path order. Maps merge by key, resources merge by name,
and explicitly specified fields in later files override earlier files.
List-valued fields such as webhook targets, events, and replication filters are
replaced as a unit, as is the replication trigger. This directly supports a
reusable base configuration, recommended defaults, required settings, and small
environment-specific overlays while the resolved document remains a regular
`HarborConfiguration`.

Apply is additive. It creates missing resources and updates changed managed
fields, but it never deletes resources omitted from the file.

## Secrets

Exports never contain credentials or webhook authentication headers. A file
that needs to create or rotate those values can reference environment variables:

```yaml
spec:
  registries:
    - name: private
      type: harbor
      url: https://registry.example.com
      credential:
        type: basic
        accessKeyFrom:
          env: HARBOR_REGISTRY_USERNAME
        accessSecretFrom:
          env: HARBOR_REGISTRY_PASSWORD
```

Existing credentials remain unchanged when `credential` is omitted.

This format is not a database or artifact backup and does not include Harbor
installation configuration, users, robot secrets, repositories, or artifacts.
