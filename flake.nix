{
  description = "go-output — Reusable Go library for CLI output formatting across 16 formats with NOM-style progress visualization";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    systems.url = "github:nix-systems/default";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    git-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.treefmt-nix.flakeModule
        inputs.git-hooks.flakeModule
      ];

      systems = import inputs.systems;

      perSystem =
        {
          config,
          pkgs,
          ...
        }:
        let
          go = pkgs.go_1_26;
          modules = [
            "."
            "bdd"
            "delimited"
            "d2"
            "enum"
            "envdetect"
            "escape"
            "examples"
            "graph"
            "integration"
            "markdown"
            "markup"
            "nom"
            "plantuml"
            "serialization"
            "table"
            "testhelpers"
            "testhelpers/graphtest"
            "tree"
            "tui"
          ];

          runForModules =
            action:
            pkgs.writeShellApplication {
              name = "go-${action}";
              runtimeInputs = [ go ];
              text = ''
                set -euo pipefail
                for mod in ${pkgs.lib.concatStringsSep " " modules}; do
                  echo ":: $mod :: go ${action} ./..."
                  ( cd "$mod" && go ${action} ./... )
                done
              '';
            };
        in
        {
          treefmt.config = {
            projectRootFile = "go.mod";

            programs = {
              nixfmt.enable = true;
              deadnix.enable = true;
              statix.enable = true;
            };
          };

          pre-commit.settings = {
            hooks = {
              treefmt.enable = true;
            };
          };

          checks.format = config.treefmt.build.check inputs.self;

          packages.default =
            pkgs.runCommand "go-output"
              {
                meta = with pkgs.lib; {
                  description = "Reusable Go library for CLI output formatting across 16 formats with NOM-style progress visualization";
                  homepage = "https://github.com/larsartmann/go-output";
                  license = licenses.mit;
                  platforms = platforms.unix;
                };
              }
              ''
                echo "go-output library — use 'go get github.com/larsartmann/go-output' to install" > $out
              '';

          devShells = {
            default = pkgs.mkShellNoCC {
              name = "go-output";

              packages = builtins.attrValues {
                inherit go;
                inherit (pkgs) golangci-lint gopls;
              };

              GOWORK = "off";

              shellHook = config.pre-commit.shellHook;
            };

            ci = pkgs.mkShellNoCC {
              name = "go-output-ci";

              packages = builtins.attrValues {
                inherit go;
                inherit (pkgs) golangci-lint;
              };

              GOWORK = "off";
            };
          };

          apps = {
            test = {
              type = "app";
              program = runForModules "test";
            };

            test-race = {
              type = "app";
              program = pkgs.writeShellApplication {
                name = "go-test-race";
                runtimeInputs = [ go ];
                text = ''
                  set -euo pipefail
                  for mod in nom tui; do
                    echo ":: $mod :: go test -race -count=1 ./..."
                    ( cd "$mod" && go test -race -count=1 ./... )
                  done
                '';
              };
            };

            build = {
              type = "app";
              program = runForModules "build";
            };

            lint = {
              type = "app";
              program = pkgs.writeShellApplication {
                name = "go-lint";
                runtimeInputs = [
                  go
                  pkgs.golangci-lint
                ];
                text = ''
                  set -euo pipefail
                  for mod in ${pkgs.lib.concatStringsSep " " modules}; do
                    echo ":: $mod :: golangci-lint run ./..."
                    ( cd "$mod" && golangci-lint run ./... )
                  done
                '';
              };
            };

            tidy = {
              type = "app";
              program = pkgs.writeShellApplication {
                name = "go-mod-tidy";
                runtimeInputs = [ go ];
                text = ''
                  set -euo pipefail
                  for mod in ${pkgs.lib.concatStringsSep " " modules}; do
                    echo ":: $mod :: go mod tidy"
                    ( cd "$mod" && go mod tidy )
                  done
                '';
              };
            };

            govulncheck = {
              type = "app";
              program = pkgs.writeShellApplication {
                name = "go-govulncheck";
                runtimeInputs = [
                  go
                  pkgs.govulncheck
                ];
                text = ''
                  set -euo pipefail
                  for mod in ${pkgs.lib.concatStringsSep " " modules}; do
                    echo ":: $mod :: govulncheck ./..."
                    ( cd "$mod" && govulncheck ./... )
                  done
                '';
              };
            };

            setup-workspace = {
              type = "app";
              program = pkgs.writeShellApplication {
                name = "setup-workspace";
                text = ''
                  if [ -f go.work ]; then
                    echo "go.work already exists, skipping"
                    exit 0
                  fi
                  if [ ! -f go.work.example ]; then
                    echo "ERROR: go.work.example not found" >&2
                    exit 1
                  fi
                  cp go.work.example go.work
                  echo "Generated go.work from go.work.example"
                '';
              };
            };
          };
        };
    };
}
