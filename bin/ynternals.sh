#!/usr/bin/env bash
set -euo pipefail

KEYFILE="/etc/nixos/symmetric.key"
ENC_FILE="$1"

if [ -z "$ENC_FILE" ]; then
    echo "Usage: ynternals <file.json.enc>"
    exit 1
fi

if [ ! -f "$KEYFILE" ]; then
    echo "ERROR: key not found at $KEYFILE"
    echo "generate one with: openssl rand -base64 32 > $KEYFILE"
    exit 1
fi

TMP_FILE=$(mktemp /dev/shm/ynternals.XXXXXX.json)
chmod 600 "$TMP_FILE"

trap 'rm -f "$TMP_FILE"' EXIT

if [ -f "$ENC_FILE" ]; then
    openssl enc -d -aes-256-cbc -pbkdf2 -in "$ENC_FILE" -out "$TMP_FILE" -pass "file:$KEYFILE" 2>/dev/null || {
        echo "failed to decrypt. wrong key or corrupted file."
        exit 1
    }
else
    echo "{}" > "$TMP_FILE"
fi

${EDITOR:-nano} "$TMP_FILE"

openssl enc -aes-256-cbc -pbkdf2 -salt -in "$TMP_FILE" -out "$ENC_FILE" -pass "file:$KEYFILE"

echo "secrets successfully saved to $ENC_FILE"
