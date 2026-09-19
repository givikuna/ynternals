{
  config,
  lib,
  pkgs,
  ...
}:
with lib;
let
  cfg = config.ynternals;
in
{
  options.ynternals = {
    enable = mkEnableOption "Ynternals symmetric secrets decryption";

    encrypted-file = mkOption {
      type = types.path;
      description = "path to the encrypted secrets.json.enc in your repo";
    };

    key-file = mkOption {
      type = types.str;
      default = "/etc/nixos/symmetric.key";
      description = "absolute path tot eh symmetric key on the live system. Do NOT make this a nix path, use a string, for security";
    };
  };

  config = mkIf cgf.enable {
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
        jq
        coreutils
      ];

      script = ''
        ${pkgs.bash}/bin/bash ${./scripts/decrypt.sh} "${cfg.encrypted-file}" "${cfg.key-file}"
      '';
    };
  };
}
