# Local Harbor for end-to-end tests

A disposable Harbor instance built from the upstream `goharbor/*` images, used by
the tests in [`test/e2e`](../e2e). It replaces the old dependency on
`demo.goharbor.io`, so the suite no longer breaks when a public demo server is
reset, rate-limited or offline.

## Quick start

```bash
./test/harbor/harbor.sh up          # start and block until the API answers
go test -tags e2e ./test/e2e/...
./test/harbor/harbor.sh reset       # stop and delete the volumes
```

Harbor is then at <http://localhost:8080>, log in with `admin` / `Harbor12345`.

## Requirements

- Docker with Compose v2, or Podman with `docker compose` / `podman-compose`
- ~2 GB of RAM and ~1 GB of disk for the images

Everything else is in this directory. `harbor.sh up` generates the token signing
key on first boot inside an init container, so no host tooling beyond a
container runtime and `curl` is needed.

## Podman

Nothing special: `harbor.sh` detects `docker compose`, `podman-compose` or
`docker-compose`, in that order. Rootless Podman works because the stack only
binds port 8080. Force a particular command with `COMPOSE`:

```bash
COMPOSE="podman-compose" ./test/harbor/harbor.sh up
```

## Configuration

| Variable             | Default          | Purpose                             |
| -------------------- | ---------------- | ----------------------------------- |
| `HARBOR_VERSION`     | `v2.15.2`        | Tag for every `goharbor/*` image    |
| `HARBOR_PORT`        | `8080`           | Host port the proxy binds           |
| `HARBOR_URL`         | `http://localhost:8080` | External URL Harbor advertises |
| `HARBOR_PASSWORD`    | `Harbor12345`    | Password for the `admin` account    |
| `LOG_LEVEL`          | `info`           | Log level for core and jobservice   |

The same `HARBOR_URL`, `HARBOR_USERNAME` and `HARBOR_PASSWORD` variables are read
by the e2e tests, so pointing them at an existing Harbor instead needs no changes
here:

```bash
HARBOR_URL=https://harbor.example.com HARBOR_USERNAME=admin HARBOR_PASSWORD=... \
  go test -tags e2e ./test/e2e/...
```

## What runs

| Service       | Role                                                    |
| ------------- | ------------------------------------------------------- |
| `proxy`       | nginx front door on port 8080, routes to core and portal |
| `core`        | the API the CLI talks to                                |
| `portal`      | the web UI                                              |
| `jobservice`  | async jobs (GC, replication, retention)                 |
| `registry`    | `docker/distribution`, blob and manifest storage        |
| `registryctl` | storage operations used by GC                           |
| `postgresql`  | Harbor's database                                       |
| `redis`       | cache and job queue                                     |

Trivy, the exporter and the syslog collector from a full Harbor deployment are
left out — the CLI does not need them and they roughly double the startup time.

## Security

The credentials and secrets in `docker-compose.yaml` are fixed, committed and
insecure by design: this is a throwaway instance bound to localhost over plain
HTTP. Do not reuse this configuration anywhere else.
