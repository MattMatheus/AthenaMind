#!/usr/bin/env bash
set -euo pipefail

issued_at="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
nonce="$(LC_ALL=C tr -dc 'A-Z0-9' </dev/urandom | head -c 20)"

echo "security_nonce=$nonce"
echo "issued_at_utc=$issued_at"
echo "window_minutes=${ATHENA_SECURITY_GATE_WINDOW_MINUTES:-15}"
echo
echo "Set these before security-affecting launch:"
echo "  export ATHENA_SECURITY_CHANGE_NONCE=$nonce"
echo "  export ATHENA_SECURITY_NONCE_ISSUED_AT_UTC=$issued_at"
