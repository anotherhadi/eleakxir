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
    vendorHash = "sha256-NDY3T3FhQ2iXJr3v3sxTX9taVTU9LPCLd/emWukHZcs=";

    buildInputs = [
      pkgs.duckdb
      pkgs.arrow-cpp
    ];
  };
in {
  package = package;
}
