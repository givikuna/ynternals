#!/usr/bin/env bash

set -euo pipefail

SECRETS_FILE="$1"
KEY_FILE="$2"

if [ ! -f "$SECRETS_FILE" ] || [ ! -f "$KEY_FILE" ]; then
    echo "secrets or key file missing."
    exit 0
fi

while IFS= read -r key && IFS= read -r encValue; do
    [ -z "$key" ] && continue

    if [[ ! "$key" =~ ^[A-Za-z0-9._-]+$ ]] || [[ "$key" == *".."* ]]; then
        echo "invalid secret name: $key" >&2
        exit 1
    fi

    printf '%s\n' "$encValue" | openssl enc -d -aes-256-cbc -pbkdf2 -salt -a -pass "file:$KEY_FILE" > "/run/ynternals/$key"

done < <(jq -r 'to_entries[] | .key, .value' "$SECRETS_FILE")
