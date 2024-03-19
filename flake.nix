{
	inputs = {
		nixpkgs.url = "github:nixos/nixpkgs?ref=85f1ba3e51676fa8cc604a3d863d729026a6b8eb";
		flake-utils.url = "github:numtide/flake-utils";
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
		{ self, nixpkgs, gomod2nix, npmlock2nix, flake-utils }:

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
				packages.default = import ./. {
					inherit pkgs;
					src = self;
				};
			}
		);
}
