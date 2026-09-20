#!/usr/bin/env bash
set -euo pipefail

SECRETS_FILE="$1"
KEY_FILE="$2"

if [ ! -f "$SECRETS_FILE" ] || [ ! -f "$KEY_FILE" ]; then
    echo "secrets or key file missing."
    exit 0
fi

mkdir -p /run/ynternals
chmod 750 /run/ynternals
chown root:wheel /run/ynternals

while IFS= read -r key && IFS= read -r encValue; do
    [ -z "$key" ] && continue

    echo "$encValue" | openssl enc -d -aes-256-cbc -pbkdf2 -salt -a -pass "file:$KEY_FILE" > "/run/ynternals/$key"

    chmod 440 "/run/ynternals/$key"
    chown root:wheel "/run/ynternals/$key"

done < <(jq -r 'to_entries[] | .key, .value' "$SECRETS_FILE")
