#!/usr/bin/env bash
# Generate Harbor configuration files for local E2E testing.
# This script creates the minimal config files needed by Harbor core.

set -euo pipefail

CONFIG_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/config/core"

mkdir -p "${CONFIG_DIR}"

# --- app.conf ---
cat > "${CONFIG_DIR}/app.conf" <<'EOF'
appname = Harbor
runmode = prod
enablegzip = true

[prod]
httpport = 8080
EOF

# --- Core encryption key (16-byte hex → 32 hex chars) ---
# This key encrypts DB passwords and secrets stored in Harbor's DB.
cat > "${CONFIG_DIR}/key" <<'EOF'
deadbeefcafebabedeadbeefcafebabe
EOF

# --- RSA private key for token signing ---
# Generate a fresh key each time so tokens are valid across restarts.
openssl genrsa -out "${CONFIG_DIR}/private_key.pem" 2048 2>/dev/null

echo "Harbor E2E config generated in ${CONFIG_DIR}"
