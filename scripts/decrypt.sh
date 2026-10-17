#!/usr/bin/env bash

set -euo pipefail

SECRETS_FILE="$1"
KEY_FILE="$2"

if [ ! -f "$SECRETS_FILE" ] || [ ! -f "$KEY_FILE" ]; then
    echo "secrets or key file missing."
    exit 0
fi

if ! jq . >/dev/null 2>&1 < "$SECRETS_FILE"; then
    echo "ERROR: $SECRETS_FILE is not valid JSON" >&2
    exit 1
fi

if ! jq -e 'type == "object"' >/dev/null 2>&1 < "$SECRETS_FILE"; then
    echo "ERROR: $SECRETS_FILE must be a JSON object" >&2
    exit 1
fi

if ! jq -e '[.[] | select(type != "string")] | length == 0' >/dev/null 2>&1 < "$SECRETS_FILE"; then
    echo "ERROR: all values in $SECRETS_FILE must be strings" >&2
    exit 1
fi

while IFS= read -r key && IFS= read -r encValue; do
    [ -z "$key" ] && continue

    if [[ ! "$key" =~ ^[A-Za-z0-9._-]+$ ]] || [[ "$key" == *".."* ]]; then
        echo "invalid secret name: $key" >&2
        exit 1
    fi

    printf '%s\n' "$encValue" | openssl enc -d -aes-256-cbc -pbkdf2 -salt -a -pass "file:$KEY_FILE" > "/run/ynternals/$key"

done < <(jq -r 'to_entries[] | .key, .value' "$SECRETS_FILE")
