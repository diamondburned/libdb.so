{
	inputs = {
		nixpkgs.url = "github:nixos/nixpkgs?ref=85f1ba3e51676fa8cc604a3d863d729026a6b8eb";
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
	};

	outputs =
		{ self, nixpkgs, gomod2nix, npmlock2nix, flake-utils, flake-compat }:

		flake-utils.lib.eachDefaultSystem (system:
			let
				overlays = [
					(self: super: {
						go = super.go_1_21;
						buildGoModule = super.buildGo121Module;
					})
					(self: super: {
						npmlock2nix = super.callPackage npmlock2nix { };
					})
					(gomod2nix.overlays.default)
				];

				pkgs = import nixpkgs {
					inherit system overlays;
				};

				version =
					if self ? rev then
						builtins.substring 0 7 self.rev
					else
						"dirty";
			in
			{
				devShells.default = pkgs.mkShell {
					packages = with pkgs; [
						nodejs
						go
						gopls
						jq
						tinygo
						gomod2nix.packages.${system}.default
					];

					GOOS = "js";
					GOARCH = "wasm";

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

				packages.vm = (pkgs.buildGoApplication {
					inherit version;
					pname = "libdb.so-vm-wasm";
					go = pkgs.go;
					src = self;
					modules = ./gomod2nix.toml;
					subPackages = [ "vm/cmd/vm-wasm" ];

					CGO_ENABLED = 0;
					doCheck = false; # none to run

					ldflags =
						[ "-s" "-w" ]
						++ (if version != "dirty" then [ "-X main.gitrev=${version}" ] else [ ]);

					postInstall = ''
						mv $out/bin/js_wasm/vm-wasm $out/bin/vm.wasm
						rmdir $out/bin/js_wasm
					'';
				}).overrideAttrs (old: old // {
					GOOS = "js";
					GOARCH = "wasm";
				});
			}
		);
}
