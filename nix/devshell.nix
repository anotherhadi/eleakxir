{pkgs, ...}: {
  devShell = pkgs.mkShell {
    buildInputs = with pkgs; [
      duckdb
      air
      # OSINT tools
      ghunt
      sherlock
      holehe
    ];
  };
}
