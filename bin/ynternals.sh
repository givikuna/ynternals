#!/usr/bin/env bash
set -euo pipefail

COMMAND="${1:-}"

if [ "$COMMAND" = "-v" ] || [ "$COMMAND" = "--version" ]; then
    echo "ynternals version @VERSION@"
    exit 0

elif [ "$COMMAND" = "add" ]; then
    KEYFILE="/etc/nixos/symmetric.key"

    if [ ! -f "$KEYFILE" ]; then
        echo "ERROR: key not found at $KEYFILE"
        exit 1
    fi

    echo -n "enter secret data: "
    IFS= read -rs SECRET_DATA
    echo

    printf '%s' "$SECRET_DATA" | openssl enc -aes-256-cbc -pbkdf2 -salt -a -pass "file:$KEYFILE" | tr -d '\n'
    echo
    exit 0

elif [ "$COMMAND" = "remove" ]; then
    SECRET_NAME="${2:-}"
    SECRETS_FILE="${3:-secrets.json}"

    if [ -z "$SECRET_NAME" ]; then
        echo "usage: ynternals remove <secret_name> [secrets.json]"
        exit 1
    fi

    if [ ! -f "$SECRETS_FILE" ]; then
        echo "ERROR: file not found at $SECRETS_FILE"
        exit 1
    fi

    TMP_FILE=$(mktemp)
    jq "del(.\"$SECRET_NAME\")" "$SECRETS_FILE" > "$TMP_FILE"
    cat "$TMP_FILE" > "$SECRETS_FILE"
    rm -f "$TMP_FILE"

    echo "removed '$SECRET_NAME' from $SECRETS_FILE"
    exit 0

else
    echo "usage:"
    echo "  ynternals add"
    echo "  ynternals remove <secret_name> [secrets.json path]"
    echo "  ynternals -v, --version"
    exit 1
fi
