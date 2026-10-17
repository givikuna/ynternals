{
  description = "symmetric secret manager for nix";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
      version = builtins.replaceStrings [ "\n" "\r" ] [ "" "" ] (builtins.readFile ./.version);
    in
    {
      nixosModules.default = import ./nix/module.nix;

      packages.${system}.default = pkgs.writeShellApplication {
        name = "ynternals";
        runtimeInputs = with pkgs; [
          openssl
          coreutils
          jq
        ];
        text = builtins.replaceStrings [ "@VERSION@" ] [ version ] (builtins.readFile ./bin/ynternals.sh);
      };

      apps.${system}.default = {
        type = "app";
        program = "${self.packages.${system}.default}/bin/ynternals";
      };
    };
}
