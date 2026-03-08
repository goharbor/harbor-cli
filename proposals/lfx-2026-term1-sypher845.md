# Proposal: LFX Mentorship CNCF - Harbor: Harbor CLI (2026 Term 1)

Author: [Sypher845](https://github.com/Sypher845)

Mentor: [Prasanth Baskar](https://github.com/bupd)

## Abstract

During this mentorship, I will complete the remaining Harbor CLI features. My main tasks include implementing the Vulnerability, Jobservice, and Distribution dashboards commands in the CLI. I will also improve the audit logs streaming functionality, secure the release pipeline, and enhance the overall usability of the CLI.

## Proposal

### 1. Vulnerability Command
This command provides a comprehensive overview of security risks across the registry. It helps users monitor the overall security posture of the system and quickly identify vulnerable artifacts that require immediate attention.

**Subcommands:**
* `summary`: Shows a security overview of the registry, total vulnerability counts by severity, scan coverage, and optionally the top dangerous artifacts and CVEs.
* `list`: Lists vulnerabilities across all scanned artifacts with support for filtering by severity, CVE ID, CVSS score range, package, project, and more.

### 2. Jobservice Command

This command lets Harbor admins monitor and manage background jobs like replication, garbage collection, and artifact scanning right from the CLI.

**Subcommands:**

* `pools list`: Lists all active worker pools with their concurrency, running/idle worker counts, and registered job types.
* `workers list`: Lists workers in a specific pool, showing what each worker is doing and whether it's busy or idle.
* `queue list`: Shows all job queues with their pending job count, latency, and paused status.
* `queue pause`: Pauses a queue so it stops picking up new jobs.
* `queue resume`: Resumes a paused queue.
* `queue clear`: Flushes all pending jobs from a queue.
* `job stop`: Stops a running job immediately.
* `job logs`: Prints the execution log of a job.

### 3. Distribution Command
In Harbor, Distribution refers to P2P content distribution using providers like Dragonfly and Kraken. The CLI already has a basic `instance` command with `create`, `list`, and `delete`. This covers implementing the remaining instance subcommands and adding full support for preheat policies, executions, tasks, and manual preheat triggers.

### 4. Improving audit logs streaming functionality

The CLI already has a `harbor logs` command with listing, filtering, sorting, and a `--follow` mode for real-time streaming. This covers improving the existing streaming functionality and adding support for listing all available audit log event types.

### 5. Improve and Secure the Release Pipeline

The release pipeline currently uses Dagger for cross-platform builds, SBOM generation, checksums, container image signing with cosign, and package builds (APK, DEB, RPM). This covers improving the overall pipeline reliability and strengthening the security posture of the release process.

### 6. Enhance Overall CLI Usability
Several PRs have already been merged to improve CLI usability, including bug fixes, better error handling, and improved command outputs. This will continue with further enhancements across existing commands to make the CLI more intuitive and reliable for daily use.

## Implementation

I will complete this work over the 12-week mentorship:
* **Weeks 1-2:** Implement the Vulnerability command with `list` and `summary` subcommands, including filters, TUI views, and output formatting.
* **Weeks 3-4:** Implement the Jobservice command with `pools list`, `workers list`, `queue list/pause/resume/clear`, and `job stop/logs`.
* **Weeks 5-6:** Implement the Distribution command, completing the remaining `instance` subcommands and adding preheat policy CRUD.
* **Weeks 7-8:** Add preheat execution and task management, and improve the audit log streaming functionality.
* **Weeks 9-10:** Improve and secure the release pipeline with better reliability and security hardening.
* **Weeks 11-12:** Enhance overall CLI usability, fix bugs, and clean up existing commands.

## Open Issues

* [#723](https://github.com/goharbor/harbor-cli/issues/723) : Vulnerability command

