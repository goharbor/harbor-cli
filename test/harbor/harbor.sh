#!/usr/bin/env bash
# Manage the throwaway Harbor used by the end-to-end tests.
#
#   ./harbor.sh up      start the stack and block until the API answers
#   ./harbor.sh down    stop the stack, keeping its volumes
#   ./harbor.sh reset   stop the stack and delete its volumes
#   ./harbor.sh logs    follow the logs of every service
#   ./harbor.sh wait    block until the API answers
#
# Works with either Docker or Podman: set COMPOSE to override the detected
# command, e.g. COMPOSE="podman-compose".
set -euo pipefail

cd "$(dirname "$0")"

HARBOR_PORT="${HARBOR_PORT:-8080}"
HARBOR_URL="${HARBOR_URL:-http://localhost:${HARBOR_PORT}}"
# Harbor needs about a minute from a cold start; allow generously more.
WAIT_TIMEOUT="${WAIT_TIMEOUT:-300}"

detect_compose() {
  if [ -n "${COMPOSE:-}" ]; then
    echo "$COMPOSE"
  elif docker compose version >/dev/null 2>&1; then
    echo "docker compose"
  elif command -v podman-compose >/dev/null 2>&1; then
    echo "podman-compose"
  elif command -v docker-compose >/dev/null 2>&1; then
    echo "docker-compose"
  else
    echo "no compose command found: install Docker Compose or podman-compose" >&2
    exit 1
  fi
}

COMPOSE_CMD="$(detect_compose)"

wait_for_harbor() {
  echo "waiting for Harbor at ${HARBOR_URL} (up to ${WAIT_TIMEOUT}s)"
  local deadline=$((SECONDS + WAIT_TIMEOUT))
  until curl -fsS "${HARBOR_URL}/api/v2.0/ping" >/dev/null 2>&1; do
    if [ "$SECONDS" -ge "$deadline" ]; then
      echo "Harbor did not become ready within ${WAIT_TIMEOUT}s" >&2
      $COMPOSE_CMD ps >&2 || true
      exit 1
    fi
    sleep 2
  done
  echo "Harbor is ready at ${HARBOR_URL} (admin / ${HARBOR_PASSWORD:-Harbor12345})"
}

case "${1:-up}" in
up)
  $COMPOSE_CMD up -d
  wait_for_harbor
  ;;
down)
  $COMPOSE_CMD down
  ;;
reset)
  $COMPOSE_CMD down --volumes
  ;;
logs)
  $COMPOSE_CMD logs -f
  ;;
wait)
  wait_for_harbor
  ;;
*)
  echo "usage: $0 {up|down|reset|logs|wait}" >&2
  exit 1
  ;;
esac
