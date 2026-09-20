{
  config,
  lib,
  pkgs,
  ...
}:
with lib;
let
  cfg = config.ynternals;

  secret-data =
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
      default = mapAttrs (name: _: "/run/ynternals/${name}") secret-data;
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

      environment = {
        SECRETS_HASH = builtins.hashString "sha256" (builtins.toJSON secret-data);
      };

      path = with pkgs; [
        openssl
        jq
        coreutils
      ];

      script = ''
        ${pkgs.bash}/bin/bash ${./scripts/decrypt.sh} "${cfg.file}" "${cfg.key-file}"
      '';
    };
  };
}
