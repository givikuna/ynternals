{ buildGoModule }:
let
  version = builtins.replaceStrings [ "\n" "\r" ] [ "" "" ] (builtins.readFile ../.version);
in
buildGoModule {
  pname = "ynternals";
  inherit version;

  src = ./..;

  vendorHash = null;

  subPackages = [
    "cmd/ynternals"
    "cmd/ynternals-decrypt"
  ];

  ldflags = [
    "-s"
    "-w"
    "-X main.version=${version}"
  ];
}
