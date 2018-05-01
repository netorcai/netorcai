{
  pkgs ? import (
    fetchTarball "https://github.com/NixOS/nixpkgs/archive/17.09.tar.gz") {},
}:

pkgs.stdenv.mkDerivation rec {
  name = "netorcai-test-env";
  env = pkgs.buildEnv { name = name; paths = buildInputs; };

  buildInputs = [
    #########
    # Build #
    #########
    pkgs.go

    ########
    # Test #
    ########
    # Testsuite
    pkgs.bats
    # Client
    pkgs.python36
    # Misc
    pkgs.psmisc
    pkgs.netcat-gnu
  ];
}
