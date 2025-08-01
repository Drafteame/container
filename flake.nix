{
  description = "Development tools";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem(system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            husky
            commitizen

            go
            goimports-reviser
            golangci-lint
            go-task
            go-mockery
            golines
            gci

            sd
            fd
          ];

          shellHook = ''
            export GOROOT="${pkgs.go}/share/go"
          '';
        };
      }
    );
}