# Robot Account Configuration File

This document describes the **YAML** format used to declare robot accounts, their lifetime, and the permissions they will receive inside a project. **JSON** is also supported.

---

## 1. File Schema

```yaml
name: "robot-name"        # Required: Name of the robot account
description: "..."        # Optional: Description of the robot account
duration: 90              # Required: Lifetime in days (-1 for no expiration)
kind: "project"           # Required: "project" or "system" - determines robot type
permissions:              # Required: Permission scopes
  - access:               # List of access items within this scope
    - resource: "repository"  # Either specify a single resource
      actions:
        - "pull"
        - "push"
    - resources:          # Or specify multiple resources
        - "artifact"
        - "scan"
      actions:
        - "read"
    kind: "project"       # Permission scope kind (project or system)
    namespace: "my-project"  # Project name for project scope, "/" for system scope
```

**Key rules**

| Key           | Type      | Notes                                                                                                |
| ------------- | --------- | ---------------------------------------------------------------------------------------------------- |
| `name`        | *string*  | Must be unique per project. Lower‑case letters, numbers, and dashes recommended.                     |
| `description` | *string*  | Optional but **highly encouraged** for auditability.                                                 |
| `duration`    | *integer* | Days until the robot expires. Use `‑1` for an unlimited lifetime. Values `0` or < `‑1` are rejected. |
| `kind`        | *string*  | "project" or "system" level robot.                                                                   |
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
kind: "project"
permissions:
  - access:
    - resource: "repository"
      actions: ["pull", "push"]
    kind: "project"
    namespace: "my-project"
```

### 3.2 Read‑Only Monitoring Robot (No Expiration)

```yaml
name: "read-only-robot"
description: "Read-only access for monitoring"
kind: "project"
duration: -1
permissions:
  - access:
    - resources: ["repository", "artifact", "scan"]
      actions: ["read", "list"]
    kind: "project"
    namespace: "my-project"
```

### 3.3 Project Admin Robot (Full Access)

```yaml
name: "project-admin-robot"
description: "Project administration tasks"
kind: "project"
duration: 180
permissions:
  - access: 
    - resource: "project"
      actions: ["*"]
    kind: "project"
    namespace: "my-project"
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
