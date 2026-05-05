#!/usr/bin/env bash
# Harbor CLI E2E smoke test.
# Requires a running Harbor instance. Set HARBOR_URL, HARBOR_USERNAME, HARBOR_PASSWORD env vars.
#
# Usage:
#   HARBOR_BIN=./harbor-dev HARBOR_URL=http://localhost:8080 bash test/e2e/smoke_test.sh

set -euo pipefail

# --- Configuration ---
HARBOR_BIN="${HARBOR_BIN:-./harbor-dev}"
HARBOR_URL="${HARBOR_URL:-http://localhost:8080}"
HARBOR_USERNAME="${HARBOR_USERNAME:-admin}"
HARBOR_PASSWORD="${HARBOR_PASSWORD:-Harbor12345}"

# Temp directory for CLI config (isolated from user config)
CONFIG_DIR="$(mktemp -d)"
export HARBOR_CLI_CONFIG="${CONFIG_DIR}/config.yaml"

# 32-byte base64-encoded key for password encryption in CI
export HARBOR_ENCRYPTION_KEY="${HARBOR_ENCRYPTION_KEY:-$(printf 'e2e-test-key-123456789012345678' | base64)}"

# Log prefix for better output in CI
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

PASS=0
FAIL=0

cleanup() {
    rm -rf "${CONFIG_DIR}"
}
trap cleanup EXIT

run_test() {
    local description="$1"
    shift
    printf "  %s ... " "${description}"
    if output=$("$@" 2>&1); then
        printf "${GREEN}PASS${NC}\n"
        PASS=$((PASS + 1))
        return 0
    else
        printf "${RED}FAIL${NC}\n"
        printf "    Command: %s\n" "$*"
        printf "    Output: %s\n" "${output}"
        FAIL=$((FAIL + 1))
        return 1
    fi
}

echo "=== Harbor CLI E2E Smoke Tests ==="
echo "Server:   ${HARBOR_URL}"
echo "User:     ${HARBOR_USERNAME}"
echo "Binary:   ${HARBOR_BIN}"
echo "Config:   ${CONFIG_DIR}"
echo ""

# Test 1: Login
run_test "harbor login" \
    "${HARBOR_BIN}" login "${HARBOR_URL}" \
    --username "${HARBOR_USERNAME}" \
    --password "${HARBOR_PASSWORD}"

# Test 2: Health check
run_test "harbor health" \
    "${HARBOR_BIN}" health

# Test 3: Project list (should succeed, may be empty)
run_test "harbor project list" \
    "${HARBOR_BIN}" project list

# Test 4: Create a project (--storage-limit required for non-interactive mode)
run_test "harbor project create" \
    "${HARBOR_BIN}" project create "e2e-smoke-test" --public --storage-limit "-1"

# Test 5: Project list (should include the new project)
run_test "harbor project list (verify new project)" \
    "${HARBOR_BIN}" project list --output-format json

# Test 6: Repository list within the project (should be empty, but command must succeed)
run_test "harbor repo list" \
    "${HARBOR_BIN}" repo list "e2e-smoke-test"

# Test 7: Delete the project
run_test "harbor project delete" \
    "${HARBOR_BIN}" project delete "e2e-smoke-test"

# Test 8: User list (requires admin)
run_test "harbor user list" \
    "${HARBOR_BIN}" user list

# Test 9: Create a user
TEST_USER="e2e-test-user"
TEST_EMAIL="e2e-test@harbor.local"
run_test "harbor user create" \
    "${HARBOR_BIN}" user create \
    --username "${TEST_USER}" \
    --email "${TEST_EMAIL}" \
    --realname "E2E Test User" \
    --password "TestPass123"

# Test 10: User list (should include new user)
run_test "harbor user list (verify new user)" \
    "${HARBOR_BIN}" user list --output-format json

# Test 11: Delete the user
run_test "harbor user delete" \
    "${HARBOR_BIN}" user delete "${TEST_USER}"

echo ""
echo "=== Results: ${PASS} passed, ${FAIL} failed ==="

if [ "${FAIL}" -gt 0 ]; then
    exit 1
fi
exit 0
