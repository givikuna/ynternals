#!/usr/bin/env bash
set -euo pipefail

ENC_FILE="$1"
KEY_FILE="$2"

SEC_DIR="/run/ynternals"
mkdir -p "$SEC_DIR"
chmod 700 "$SEC_DIR"

TMP_JSON=$(mktemp)

openssl enc -d -aes-256-cbc -pbkdf2 \
  -in "$ENC_FILE" \
  -out "$TMP_JSON" \
  -pass "file:$KEY_FILE"

for key in $(jq -r 'keys[]' "$TMP_JSON"); do
  jq -r --arg k "$key" '.[$k]' "$TMP_JSON" > "$SEC_DIR/$key"
  chmod 400 "$SEC_DIR/$key"
  chown root:root "$SEC_DIR/$key"
done

shred -u "$TMP_JSON"
