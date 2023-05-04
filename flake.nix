{
  # main
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    devshell = {
      url = "github:numtide/devshell";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    ...
  } @ inputs:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = with inputs; [
            devshell.overlays.default
          ];
        };
      in {
        packages.default = pkgs.buildGoModule {
          pname = "meilisearch-syncer";
          version = "0.1.0";
          src = pkgs.lib.cleanSourceWith {
            src = ./.;
            filter = path: type: let p = baseNameOf path; in !((p == "flake.nix") || (p == "flake.lock") || (p == "README.md"));
          };
          vendorSha256 = "sha256-eYXDJUbXc5OEubhspDj9gd278pxZNFaeE5Sz4hNUxt0=";
        };
        packages.docker = pkgs.dockerTools.buildImage {
          name = "meilisearch-syncer";
          tag = "latest";
          copyToRoot = pkgs.buildEnv {
            name = "image-root";
            paths = [self.packages.${system}.default];
            pathsToLink = ["/bin"];
          };
          config.Cmd = ["/bin/meilisearch-syncer"];
        };
        devShells.default = pkgs.devshell.mkShell {
          packages = with pkgs; [
            actionlint
            alejandra
            dive
            go
            go-outline
            gocode-gomod
            gopkgs
            gopls
            gotools
            hadolint
            httpie
          ];
        };
      }
    );
}
