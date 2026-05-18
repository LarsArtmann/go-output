{
  description = "go-output — Reusable Go library for CLI output formatting across 12 formats";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    flake-parts.url = "github:hercules-ci/flake-parts";
    flake-parts.inputs.nixpkgs-lib.follows = "nixpkgs";

    treefmt-nix.url = "github:numtide/treefmt-nix";
    treefmt-nix.inputs.nixpkgs.follows = "nixpkgs";

    git-hooks.url = "github:cachix/git-hooks.nix";
    git-hooks.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.treefmt-nix.flakeModule
        inputs.git-hooks.flakeModule
      ];

      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      perSystem =
        {
          config,
          pkgs,
          ...
        }:
        let
          go = pkgs.go_1_26;
        in
        {
          # ── Formatter (nix fmt) ──────────────────────────────
          treefmt.config = {
            projectRootFile = "flake.nix";

            programs = {
              # Nix
              nixfmt.enable = true;
              deadnix.enable = true;
              statix.enable = true;
            };
          };

          # ── Pre-commit hooks ─────────────────────────────────
          pre-commit.settings = {
            hooks = {
              treefmt.enable = true;
            };
          };

          # ── Dev Shell (nix develop) ──────────────────────────
          devShells.default = pkgs.mkShellNoCC {
            name = "go-output";

            packages = builtins.attrValues {
              inherit go;
              inherit (pkgs) golangci-lint gopls;
            };

            shellHook = ''
              ${config.pre-commit.shellHook}
              export GOWORK=off
            '';
          };
        };
    };
}
