{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    flake-compat.url = "https://flakehub.com/f/edolstra/flake-compat/1.tar.gz";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs = {
        nixpkgs.follows = "nixpkgs";
        flake-utils.follows = "flake-utils";
      };
    };
    npmlock2nix = {
      url = "github:nix-community/npmlock2nix";
      flake = false;
    };
    yaml-language-server-src = {
      url = "github:okybr/yaml-language-server";
      flake = false;
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      gomod2nix,
      npmlock2nix,
      flake-utils,
      flake-compat,
      yaml-language-server-src,
    }:

    flake-utils.lib.eachDefaultSystem (
      system:
      let
        overlays = [
          (self: super: {
            npmlock2nix = import npmlock2nix {
              pkgs = super;
              lib = super.lib;
            };
          })
          (gomod2nix.overlays.default)
        ];

        pkgs = import nixpkgs { inherit system overlays; };
        lib = pkgs.lib;

        go = pkgs.go_1_22;

        nodejs = pkgs.nodejs;

        tinygo = pkgs.tinygo;
        # let
        #   overrides = rec {
        #     version = "0.32.0";
        #     src = pkgs.fetchFromGitHub {
        #       owner = "tinygo-org";
        #       repo = "tinygo";
        #       rev = "v${version}";
        #       hash = "sha256-zoXruGoWitx6kietF3HKTYCtUrXp5SOrf2FEGgVPzkQ=";
        #       fetchSubmodules = true;
        #     };
        #     doCheck = false;
        #     patches = [ ];
        #   };
        # in
        # (pkgs.tinygo.overrideAttrs (_: overrides)).override {
        #   buildGoModule =
        #     args:
        #     pkgs.buildGoModule (
        #       args // overrides // { vendorHash = "sha256-rJ8AfJkIpxDkk+9Tf7ORnn7ueJB1kjJUBiLMDV5tias="; }
        #     );
        # };

        yaml-language-server = pkgs.yaml-language-server.overrideAttrs (old: rec {
          version = "1.15.0-ajv-draft-04";
          src = yaml-language-server-src;
          offlineCache = pkgs.fetchYarnDeps {
            yarnLock = "${src}/yarn.lock";
            hash = "sha256-thJ3aU52yCusfjBCD2QvLynwiM32lq0IT9WaNJjfu6E=";
          };
        });

        gopls =
          let
            GOOS = "js";
            GOARCH = "wasm";
          in
          pkgs.writeShellScriptBin "gopls" ''
            export GOOS=${GOOS}
            export GOARCH=${GOARCH}
            exec ${lib.getExe pkgs.gopls} "$@"
          '';

        version = if self ? rev then builtins.substring 0 7 self.rev else "dirty";
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            nodejs
            go
            gopls
            jq
            yq-go
            # tinygo
            oapi-codegen
            yaml-language-server
            self.formatter.${system}
            gomod2nix.packages.${system}.default
          ];

          shellHook = ''
            export PATH="$PATH:$(git rev-parse --show-toplevel)/node_modules/.bin"
          '';
        };

        packages.default = pkgs.stdenv.mkDerivation rec {
          inherit version;
          pname = "libdb.so";
          src = self;

          nativeBuildInputs = with pkgs; [
            coreutils
            bash
            jq
            nodejs
          ];

          nodeModules = pkgs.npmlock2nix.v2.node_modules {
            inherit src;
            nodejs = pkgs.nodejs;
            # mkDerivation hates us because we have a Makefile. We'll override
            # installPhase to fix that.
            installPhase = "mv node_modules $out/";
          };

          preBuild = ''
            set -x

            mkdir -p build
            cp -r ${self.packages.${system}.vm}/bin/vm.wasm build/vm.wasm

            cp -r ${nodeModules} node_modules
            chown -R $(id -u):$(id -g) node_modules
            chmod -R +w node_modules
            export PATH="$PATH:$PWD/node_modules/.bin"
            export VERSION="$version"

            set +x
          '';

          installPhase = ''
            cp -r build/dist $out
          '';
        };

        packages.vm =
          (pkgs.buildGoApplication {
            inherit version go;
            pname = "libdb.so-vm-wasm";
            src = self;
            modules = ./gomod2nix.toml;
            subPackages = [ "vm/cmd/vm-wasm" ];

            CGO_ENABLED = 0;
            doCheck = false; # none to run

            ldflags = [
              "-s"
              "-w"
            ] ++ (if version != "dirty" then [ "-X main.gitrev=${version}" ] else [ ]);

            postInstall = ''
              mv $out/bin/js_wasm/vm-wasm $out/bin/vm.wasm
              rmdir $out/bin/js_wasm
            '';
          }).overrideAttrs
            (
              old:
              old
              // {
                GOOS = "js";
                GOARCH = "wasm";
              }
            );

        packages.backend = pkgs.buildGoApplication {
          inherit version go;
          pname = "libdb.so-backend";
          src = self;
          modules = ./gomod2nix.toml;
          subPackages = [ "backend/cmd/backend" ];
        };

        formatter = pkgs.nixfmt-rfc-style;
      }
    );
}
