{
  description = "Flake for eleakxir";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
  };

  outputs = {
    systems,
    self,
    nixpkgs,
    ...
  }: let
    eachSystem = nixpkgs.lib.genAttrs (import systems);

    importBackend = system:
      import ./nix/back.nix {
        pkgs = nixpkgs.legacyPackages.${system};
        lib = nixpkgs.lib;
        inherit self;
      };

    importUtils = system:
      import ./nix/leak-utils.nix {
        pkgs = nixpkgs.legacyPackages.${system};
        lib = nixpkgs.lib;
        inherit self;
      };

    importDevShell = system:
      import ./nix/devshell.nix {
        pkgs = nixpkgs.legacyPackages.${system};
        lib = nixpkgs.lib;
        inherit self;
      };
  in {
    packages = eachSystem (system: {
      backend = (importBackend system).package;
      leak-utils = (importUtils system).package;
    });

    devShells = eachSystem (system: {
      default = (importDevShell system).devShell;
    });

    nixosModules.eleakxir = {
      config,
      lib,
      pkgs,
      ...
    }:
      (importBackend pkgs.system).nixosModule {
        inherit config lib pkgs;
        inherit self;
      };
  };
}
