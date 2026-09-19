{
  config,
  lib,
  pkgs,
  ...
}:
with lib;
let
  cfg = config.ynternals;

  secretData =
    if cfg.file != null && builtins.pathExists cfg.file then
      builtins.fromJSON (builtins.readFile cfg.file)
    else
      { };
in
{
  options.ynternals = {
    enable = mkEnableOption "Ynternals symmetric secrets decryption";

    file = mkOption {
      type = types.nullOr types.path;
      default = null;
      description = "Path to your secrets.json file";
    };

    key-file = mkOption {
      type = types.str;
      default = "/etc/nixos/symmetric.key";
      description = "Absolute path to the symmetric key on the live system. Do NOT make this a nix path, use a string, for security";
    };

    secrets = mkOption {
      type = types.attrsOf types.str;
      readOnly = true;
      description = "Attribute set mapping secret names to their decrypted file paths";
      default = mapAttrs (name: _: "/run/ynternals/${name}") secretData;
    };
  };

  config = mkIf cfg.enable {
    systemd.services.ynternals-decrypt = {
      description = "Decrypt ynternals secrets";
      wantedBy = [ "multi-user.target" ];
      before = [ "network.target" ];

      serviceConfig = {
        Type = "oneshot";
        RemainAfterExit = true;
      };

      path = with pkgs; [
        openssl
      ];

      script = ''
        set -euo pipefail
        mkdir -p /run/ynternals
        chmod 700 /run/ynternals

        ${concatStringsSep "\n" (
          mapAttrsToList (name: encValue: ''
            echo "${encValue}" | openssl enc -d -aes-256-cbc -pbkdf2 -salt -a -pass "file:${cfg.key-file}" > "/run/ynternals/${name}"
            chmod 400 "/run/ynternals/${name}"
          '') secretData
        )}
      '';
    };
  };
}
