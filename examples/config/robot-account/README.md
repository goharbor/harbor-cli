# Robot Account Configuration File

This document describes the **YAML** format used to declare robot accounts, their lifetime, and the permissions they will receive inside a project. **JSON** is also supported.

---

## 1. File Schema

```yaml
name: "robot-name"        # Required – unique name for the robot account
description: "..."        # Optional – free‑form description
duration: 90              # Required – lifetime in days (‑1 means no expiration)
project: "project-name"   # Required – target project
permissions:              # Required – at least one rule (see below)
  # Either a single resource …
  - resource: "repository"
    actions: ["pull", "push"]
  # … or multiple resources that share the same actions
  - resources: ["artifact", "scan"]
    actions: ["read"]
  # Grant every action for a single resource with '*'
  - resource: "project"
    actions: ["*"]
```

**Key rules**

| Key           | Type      | Notes                                                                                                |
| ------------- | --------- | ---------------------------------------------------------------------------------------------------- |
| `name`        | *string*  | Must be unique per project. Lower‑case letters, numbers, and dashes recommended.                     |
| `description` | *string*  | Optional but **highly encouraged** for auditability.                                                 |
| `duration`    | *integer* | Days until the robot expires. Use `‑1` for an unlimited lifetime. Values `0` or < `‑1` are rejected. |
| `project`     | *string*  | The Harbor project where the robot lives.                                                            |
| `permissions` | *list*    | One or more permission blocks, each granting one **set of actions** to one or more **resources**.    |

> \*\*Tip \*\*: Store separate YAML files per robot to keep Git history clean and roll back permission changes safely.

---

## 2. Defining Permissions

A permission block has three fields:

```yaml
- resource: "repository"     # or "resources" if you list several
  actions: ["pull", "push"]   # zero‑or‑more supported verbs
```

### 2.1 Single Resource, Specific Actions

```yaml
- resource: "repository"
  actions: ["pull", "push"]
```

### 2.2 Multiple Resources, Same Actions

```yaml
- resources: ["artifact", "scan"]
  actions: ["read"]
```

### 2.3 All Actions for a Resource

```yaml
- resource: "project"
  actions: ["*"]
```

> **Wildcard `*`** always means *all actions supported by that resource* – nothing more, nothing less.

---

## 3. End‑to‑End Examples

### 3.1 CI Pipeline Robot (Push/Pull Images)

```yaml
name: "ci-pipeline-robot"
description: "Robot account for CI/CD pipeline"
duration: 90
project: "my-project"
permissions:
  - resource: "repository"
    actions: ["pull", "push"]
```

### 3.2 Read‑Only Monitoring Robot (No Expiration)

```yaml
name: "read-only-robot"
description: "Read-only access for monitoring"
duration: -1
project: "my-project"
permissions:
  - resources: ["repository", "artifact", "scan"]
    actions: ["read", "list"]
```

### 3.3 Project Admin Robot (Full Access)

```yaml
name: "project-admin-robot"
description: "Project administration tasks"
duration: 180
project: "my-project"
permissions:
  - resource: "project"
    actions: ["*"]
```

---

## 4. Available Resource Types

* **repository** – Docker/OCI repositories
* **artifact** – Container images & other artifacts
* **tag** – Image tags
* **scan** – Vulnerability scan results
* **label** – Repository/artifact labels
* **project** – Project‑level settings & configs
* **sbom** – Software Bill of Materials documents
* **metadata** – Custom project metadata
* **member** – Project members & roles
* **robot** – Robot accounts themselves
* *…and others added in future Harbor versions*

---

## 5. Common Actions (per Resource)

```
repository : pull, push, delete, list, read
artifact   : read, list, create, delete
project    : read, update, delete
scan       : read, create, stop
```

> Check your Harbor version docs for the authoritative list; action sets can evolve.

---

## 6. Best Practices & Pitfalls

1. **Least Privilege First** – start with `read` and add only what the robot really needs.
2. **Short Expirations** – prefer a finite `duration`; rotate credentials via CI secrets managers.
3. **One Robot ≠ Many Jobs** – create separate robots per pipeline to keep scopes narrow.
4. **Version Control** – commit the YAML and review via pull requests.
5. **Avoid `*`** on critical resources unless you truly need admin‑like power.
