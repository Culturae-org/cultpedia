{
  description = "Cultpedia - a knowledge base for Culturae";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
  flake-utils.lib.eachSystem [ 
    "x86_64-linux" 
    "aarch64-linux" 
    "x86_64-darwin" 
    "aarch64-darwin" 
    ] (system:
      let
        pkgs = import nixpkgs { 
          inherit system;
          config.allowUnfree = true;
        };
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go

            # Only for code quality checks, it's just for developers not dataset maintainers
            # golangci-lint
          ];
        };
      });
}