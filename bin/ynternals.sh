#!/usr/bin/env bash
set -euo pipefail

if [ "${1:-}" != "add" ]; then
    echo "Usage: ynternals add"
    exit 1
fi

KEYFILE="/etc/nixos/symmetric.key"

if [ ! -f "$KEYFILE" ]; then
    echo "ERROR: key not found at $KEYFILE"
    exit 1
fi

echo -n "enter secret data: "
read -rs SECRET_DATA
echo

echo -n "$SECRET_DATA" | openssl enc -aes-256-cbc -pbkdf2 -salt -a -pass "file:$KEYFILE"
echo
