# Security Policy

Ynternals is currently in its Beta.
Therefore, only the latest released versions receive security fixes.

If you are on an older version, upgrade before reporting any issues.

## Reporting a Vulnerability

**Please do NOT open a public GitHub issue for undisclosed vulnerabilities.**

Use GitHub's private vulnerability reporting instead:

- Go to <https://github.com/givikuna/ynternals/security/advisories>
- Click **Report Vulnerability**
- Fill in all the details

Please include:

- The version tested on
- Description of the issue
- How to reproduce vulnerability
- Any suggested fixes

## Disclosure Policy

We follow coordinated disclosure. Once a fix is available, we will publish
a GitHub Security Advisory describing the issue, the affected versions,
and the remediation, and credit the reporter (unless anonymity is
requested). Please give us a reasonable window to ship a fix before any
public disclosure.

## Current Security Model

### Cryptography

- **Cipher:** AES-256-CBC.
- **Key derivation:** PBKDF2-HMAC-SHA256 with a per-secret, randomly
  generated 8-byte salt (`crypto/rand`).
- **Container format:** OpenSSL-compatible `Salted__` header followed by
  the salt and ciphertext, Base64-encoded. Ciphertexts produced by
  Ynternals can be decrypted with
  `openssl enc -d -aes-256-cbc -pbkdf2` using the same key file, and vice
  versa.
- **Padding:** PKCS#7, validated on decrypt.

### Key Handling

- The symmetric key lives in a file on the live system (default:
  `/etc/nixos/symmetric.key`, overridable via `YNTERNALS_KEYFILE`).
- The key file path is passed to the Nix module as a **string**, never as
  a Nix path, so the key is never copied into the Nix store.
- Only ciphertext is stored in `secrets.json`, and therefore only
  ciphertext ever enters the Nix store.
- Encryption is performed on the operator's machine by `ynternals add`;
  the plaintext is read from the terminal with echo disabled and is never
  written to disk.

### Runtime Decryption

- `ynternals-decrypt` runs as a **root** `oneshot` systemd service.
- Decrypted secrets are written to `/run/ynternals` (tmpfs, i.e.
  memory-backed, not persisted to disk).
- The runtime directory is `0750 root:wheel`; secret files are written
  with mode `0666` and a `UMask` of `0226`, resulting in `0640
root:wheel` files. Only root and members of `wheel` can read them.
- `ynternals peek` requires `EUID 0` and is intended for development only.

### Input Validation

- Secret names must match `^[A-Za-z0-9._-]+$` and may not contain `..`.
  This is enforced both in `ynternals add` and in the decrypt service
  before any file is written, providing path-traversal protection.
- `secrets.json` must be a flat JSON object whose values are all strings.
  Any other shape is rejected by `ynternals validate` and by the decrypt
  service.
- All file writes go through a staging path (`WriteAtomicTemp`) rather
  than being done in place.
