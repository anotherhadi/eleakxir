{
  pkgs,
  lib,
  self,
}: let
  name = "eleakxir";

  package = pkgs.buildGoModule {
    pname = name;
    version = "0.1.0";
    src = ../back;
    vendorHash = "";

    buildInputs = [
      pkgs.duckdb
      pkgs.arrow-cpp
    ];
  };
in {
  package = package;

  nixosModule = {config, ...}: let
    cfg = config.services."${name}";
  in {
    options.services."${name}" = {
      enable = lib.mkEnableOption "Enable the ${name} service";
      user = lib.mkOption {
        type = lib.types.str;
        default = name;
        description = "User to run the ${name} service as";
      };
      group = lib.mkOption {
        type = lib.types.str;
        default = name;
        description = "Group to run the ${name} service as";
      };
      port = lib.mkOption {
        type = lib.types.port;
        default = 9198;
        description = "Port for the ${name} service";
      };
      folders = lib.mkOption {
        type = lib.types.listOf lib.types.str;
        default = [];
        description = "Folders to monitor for parquet files";
      };
      cacheFolder = lib.mkOption {
        type = lib.types.str;
        default = "";
        description = "Cache folder";
      };
      limit = lib.mkOption {
        type = lib.types.int;
        default = 200;
        description = "Limit of results to return";
      };
      password = lib.mkOption {
        type = lib.types.str;
        default = "";
        description = "Password for auth (empty means no auth)";
      };
      debug = lib.mkOption {
        type = lib.types.bool;
        default = false;
        description = "Debug mode";
      };
      maxCacheDuration = lib.mkOption {
        type = lib.types.str;
        default = "24h";
        description = "Max result cache duration (30m, 2h, 1d)";
      };
      reloadDataleaksInterval = lib.mkOption {
        type = lib.types.str;
        default = "1h";
        description = "Interval to reload dataleaks (30m, 2h, 1d)";
      };
      minimumQueryLength = lib.mkOption {
        type = lib.types.int;
        default = 3;
        description = "Minimum query length";
      };
      baseColumns = lib.mkOption {
        type = lib.types.listOf lib.types.str;
        default = [];
        description = "Base columns are used when the column searched is 'all'";
      };
      githubRecon = lib.mkOption {
        type = lib.types.bool;
        default = true;
        description = "Activate the github-recon OSINT module";
      };
      githubToken = lib.mkOption {
        type = lib.types.str;
        default = "";
        description = "GitHub token to use for Github recon";
      };
      githubDeepMode = lib.mkOption {
        type = lib.types.bool;
        default = false;
        description = "Activate the github-recon deep mode";
      };
    };

    config = lib.mkIf cfg.enable {
      users.users."${cfg.user}" = {
        isSystemUser = true;
        group = cfg.group;
      };
      users.groups."${cfg.group}" = {};

      systemd.services."${name}" = {
        description = "${name} service";
        after = ["network.target"];
        wantedBy = ["multi-user.target"];
        serviceConfig = {
          ExecStart = "${self.packages.${pkgs.system}.backend}/bin/cmd";
          Restart = "always";
          User = cfg.user;
          Group = cfg.group;
          StateDirectory = name;
          ReadWritePaths = ["/var/lib/${name}"];
          WorkingDirectory = "/var/lib/${name}";

          Environment = [
            "PORT=${toString cfg.port}"
            "DATALEAKS_FOLDERS=${lib.strings.concatStringsSep "," cfg.folders}"
            "DATALEAKS_CACHE_FOLDER=${cfg.cacheFolder}"
            "LIMIT=${toString cfg.limit}"
            "PASSWORD=${toString cfg.password}"
            "DEBUG=${
              if cfg.debug
              then "true"
              else "false"
            }"
            "MAX_CACHE_DURATION=${cfg.maxCacheDuration}"
            "RELOAD_DATALEAKS_INTERVAL=${cfg.reloadDataleaksInterval}"
            "MINIMUM_QUERY_LENGTH=${toString cfg.minimumQueryLength}"
            "BASE_COLUMNS=${lib.strings.concatStringsSep "," cfg.baseColumns}"
            "GITHUB_RECON=${
              if cfg.githubRecon
              then "true"
              else "false"
            }"
            "GITHUB_TOKEN=${cfg.githubToken}"
            "GITHUB_DEEP_MODE=${
              if cfg.githubDeepMode
              then "true"
              else "false"
            }"
          ];
        };
      };
    };
  };
}
