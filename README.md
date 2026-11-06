# Ynternals

Declarative, symmetric-key management system for NixOS hosts.

Ynternals allows you to encrypt secrets in a symmetric manner (using <a href="https://en.wikipedia.org/wiki/Advanced_Encryption_Standard">AES-256</a>).
Making setting up and managing huge NixOS fleets easier and more comfortable.

## Philosophy

Ynternals (as well as the other `y` packages like <a href="https://github.com/givikuna/ysyldur">Ysyldur</a> and <a href="https://github.com/givikuna/yndrium">Yndrium</a>) follow a philosophy of simplicity and hackability.

Ynternals does NOT expose your keys to the Nix store.

So, instead of relying on complex secrets management techniques, backup techniques, and so on, gives you a far simpler framework to set up fleets faster, albeit with some downsides.

## Security

Ynternals is a secure tool.

It guarantees:

- Zero nix-store leakage
- Memory-only decryption
- Strong encryption
- Path traversal protection

### Note on Security

Ynternals is secure, and safe to use in many scenarios.

But as of right now, it is under heavy development, subject to occasional breaking changes, possibly unknown vulnerabilities, and still undecided design decisions.

For use in personal configurations, Ynternals is safe to use as a way to encrypt your public nix configuration data.

Currently these are the known issues:

- Keyfile exposure on multi-user systems
- Service dependency ordering
- Silent failures on boot if key file or secrets file is missing
- Lack of granular permissions

They are documented in more depth on GitHub issues.

## NixOS Module

This flake exposes a module with the following options:

`ynternals.enable` (`boolean`) - enable ynternals.
`ynternals.file` (`path`, optional) - set the `secrets.json` file path.
`ynternals.key-file` (`string`) - the absolute path to where the master key is located.
`ynternals.secrets` (`readOnly attrset`) - dynamically generated mapping of secret names to their decrypted file paths.

## Ynternals CLI

The Ynternals CLI is designed with the purpose of making Ynternals easy to use on your hosts.

### `ynternals add <secret_name> [secrets.json]`

If `secrets.json` does not exist, this initializes it with `{}`.
Afterwards, it prompts the user for plaintext secret, and then stores it the secrets file.

### `ynternals remove <secret_name> [secrets.json]`

Removes a secret from the system.

### `ynternals validate [secrets.json]`

manually runs JSOn validation checks to check if `secrets.json` is in the correct format.

### `ynternals list [secrets.json]`

Lists all valid secrets managed by Ynternals.

### `ynternals peek <secret_name>`

Requires `sudo`, and `cat`s the secret to the user.
This is to be used in development.

Do NOT use `ynternals peek` in production.

### `ynternals -v` \/ `ynternals --version`

Prints out the current Ynternals version.

## Installation

Ynternals is packaged as a nix flake.

As of right now, there is no other supported way to use Ynternals besides as a flake.

An example installation:

```Nix
{
    inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
        ynternals.url = "github:givikuna/ynternals";
    };

    outputs = { self, nixpkgs, ynternals, ... }: {
        nixosConfigurations.HOSTNAME = nixpkgs.lib.nixosSystem {
            system = "x86_64-linux"; # or any other system
            modules = [
                ynternals.nixosModules.default;
            ];
        };
    }
}
```

To add the package:
In any file imported into your `configuration.nix` write:

```Nix
{ pkgs, inputs, ... }: {
    ynternals = {
        enable = true;
        file   = ./PATH/TO/secrets.json;
        # rest of your configuration
    };

    environment.systemPackages = [
        inputs.ynternals.packages.${pkgs.system}.default
    ];
}
```
