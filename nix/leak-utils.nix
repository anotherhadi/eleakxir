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
    vendorHash = "sha256-rTfbXCiwv/+tVXZmgztt088Zhz0OQaVTfvxXVzw4o4Q=";

    buildInputs = [
      pkgs.duckdb
      pkgs.arrow-cpp
    ];
  };
in {
  package = package;
}
