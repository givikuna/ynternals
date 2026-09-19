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
    in
    {
      nixosModules.default = import ./module.nix;

      packages.${system}.default = pkgs.writeShellApplication {
        name = "ynternals";
        runtimeInputs = with pkgs; [
          openssl
          coreutils
        ];
        text = builtins.readFile ./bin/ynternals.sh;
      };

      apps.${system}.default = {
        type = "app";
        program = "${self.packages.${system}.default}/bin/ynternals";
      };
    };
}
