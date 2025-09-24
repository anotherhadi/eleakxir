{
  pkgs,
  lib,
  self,
}: let
  name = "leak-utils";

  package = pkgs.buildGoModule {
    pname = name;
    version = "0.1.0";
    src = ../leak-utils;
    vendorHash = "sha256-qgDqmEgL7B8FvoKNwLG0buLmg9Yt54cyWwmXBifgr/g=";

    buildInputs = [
      pkgs.duckdb
      pkgs.arrow-cpp
    ];
  };
in {
  package = package;
}
