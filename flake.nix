{
  description = "go-output — Reusable Go library for CLI output formatting across 16 formats with NOM-style progress visualization";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };

    # go-standard (flakeModules.go-standard) composes treefmt-nix + systems
    # internally — this repo no longer declares them as direct inputs.
    go-nix-helpers = {
      url = "github:LarsArtmann/go-nix-helpers";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    git-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{
      flake-parts,
      ...
    }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.go-nix-helpers.flakeModules.go-standard
        inputs.git-hooks.flakeModule
      ];

      go-standard = {
        pname = "go-output";
        description = "Reusable Go library for CLI output formatting across 16 formats with NOM-style progress visualization";
        goPkgAttr = "go_1_27";
        goExperiment = "jsonv2";
        vendorHash = "sha256-R5O4GawlGbhk39eTMPRMgNUou69i7r8FWIJ3LeYEcbk=";
        # Library repo: the package exists so the ROOT module actually
        # compiles in CI (the old placeholder never built anything), but
        # tests run via apps.test across all 19 modules, not in checkPhase.
        enableCheck = false;
        # Old flake had no test check — keep the surface exactly additive
        # otherwise. (The generated overlay is a new, useful consumer
        # surface: exposes the built library package.)
        enableTestCheck = false;
        # Old treefmt ran nixfmt/deadnix/statix only — no Go formatters.
        enableGofumpt = false;
        enableGoimports = false;
        enableNixfmt = true;
        extraMeta = {
          homepage = "https://github.com/larsartmann/go-output";
          platforms = inputs.nixpkgs.lib.platforms.unix;
        };
      };

      perSystem =
        {
          config,
          pkgs,
          ...
        }:
        let
          inherit (pkgs) lib;
          go = pkgs.go_1_27;
          modules = [
            "."
            "bdd"
            "daghtml"
            "delimited"
            "d2"
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
                export GOEXPERIMENT=jsonv2
                for mod in ${lib.concatStringsSep " " modules}; do
                  echo ":: $mod :: go ${action} ./..."
                  ( cd "$mod" && go ${action} ./... )
                done
              '';
            };
        in
        {
          treefmt = {
            projectRootFile = "go.mod";
            programs = {
              deadnix.enable = true;
              statix.enable = true;
            };
          };

          pre-commit.settings = {
            hooks = {
              treefmt.enable = true;
            };
          };

          # The old hand-rolled shell, now with GOTOOLCHAIN hardening; kept
          # via mkForce so the pre-commit hook installation composes.
          devShells.default = lib.mkForce (
            pkgs.mkShellNoCC {
              name = "go-output";

              packages = builtins.attrValues {
                inherit go;
                inherit (pkgs) golangci-lint gopls govulncheck;
              };

              GOWORK = "off";
              GOTOOLCHAIN = "local";
              GOEXPERIMENT = "jsonv2";

              shellHook = config.pre-commit.shellHook;
            }
          );

          apps = {
            test = {
              type = "app";
              program = runForModules "test";
            };

            test-race = {
              type = "app";
              program = pkgs.lib.getExe (
                pkgs.writeShellApplication {
                  name = "go-test-race";
                  runtimeInputs = [ go ];
                  text = ''
                    set -euo pipefail
                    export GOEXPERIMENT=jsonv2
                    for mod in nom tui; do
                      echo ":: $mod :: go test -race -count=1 ./..."
                      ( cd "$mod" && go test -race -count=1 ./... )
                    done
                  '';
                }
              );
            };

            test-race-all = {
              type = "app";
              program = pkgs.lib.getExe (
                pkgs.writeShellApplication {
                  name = "go-test-race-all";
                  runtimeInputs = [ go ];
                  text = ''
                    set -euo pipefail
                    export GOEXPERIMENT=jsonv2
                    for mod in ${lib.concatStringsSep " " modules}; do
                      echo ":: $mod :: go test -race -count=1 ./..."
                      ( cd "$mod" && go test -race -count=1 ./... )
                    done
                  '';
                }
              );
            };

            build = {
              type = "app";
              program = runForModules "build";
            };

            lint = {
              type = "app";
              program = pkgs.lib.getExe (
                pkgs.writeShellApplication {
                  name = "go-lint";
                  runtimeInputs = [
                    go
                    pkgs.golangci-lint
                  ];
                  text = ''
                    set -euo pipefail
                    export GOEXPERIMENT=jsonv2
                    for mod in ${lib.concatStringsSep " " modules}; do
                      echo ":: $mod :: golangci-lint run ./..."
                      ( cd "$mod" && golangci-lint run ./... )
                    done
                  '';
                }
              );
            };

            tidy = {
              type = "app";
              program = pkgs.lib.getExe (
                pkgs.writeShellApplication {
                  name = "go-mod-tidy";
                  runtimeInputs = [ go ];
                  text = ''
                    set -euo pipefail
                    export GOEXPERIMENT=jsonv2
                    for mod in ${lib.concatStringsSep " " modules}; do
                      echo ":: $mod :: go mod tidy"
                      ( cd "$mod" && go mod tidy )
                    done
                  '';
                }
              );
            };

            govulncheck = {
              type = "app";
              program = pkgs.lib.getExe (
                pkgs.writeShellApplication {
                  name = "go-govulncheck";
                  runtimeInputs = [
                    go
                    pkgs.govulncheck
                  ];
                  text = ''
                    set -euo pipefail
                    export GOEXPERIMENT=jsonv2
                    for mod in ${lib.concatStringsSep " " modules}; do
                      echo ":: $mod :: govulncheck ./..."
                      ( cd "$mod" && govulncheck ./... )
                    done
                  '';
                }
              );
            };

            website-build = {
              type = "app";
              program = pkgs.lib.getExe (
                pkgs.writeShellApplication {
                  name = "website-build";
                  runtimeInputs = [
                    pkgs.nodejs
                    pkgs.pnpm
                  ];
                  text = ''
                    set -euo pipefail
                    export CI=true
                    cd website
                    pnpm install --frozen-lockfile
                    pnpm run verify
                  '';
                }
              );
            };

            website-deploy = {
              type = "app";
              program = pkgs.lib.getExe (
                pkgs.writeShellApplication {
                  name = "website-deploy";
                  runtimeInputs = [
                    pkgs.nodejs
                    pkgs.pnpm
                    pkgs.firebase-tools
                  ];
                  text = ''
                    set -euo pipefail
                    export CI=true
                    cd website
                    pnpm install --frozen-lockfile
                    pnpm run verify
                    firebase deploy --only hosting:go-output --project lars-software
                  '';
                }
              );
            };

            setup-workspace = {
              type = "app";
              program = pkgs.lib.getExe (
                pkgs.writeShellApplication {
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
                }
              );
            };
          };
        };
    };
}
