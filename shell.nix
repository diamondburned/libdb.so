{ pkgs ? import <nixpkgs> {} }:

let sources = import ./nix/sources.nix;
	ourPkgs = import sources.nixpkgs {
		overlays = [
			(self: super: {
				go = super.go_1_20;
				buildGoModule = super.buildGo120Module;
			})
		];
	};

	lib = pkgs.lib;

	tinygo = ourPkgs.callPackage ./nix/packages/tinygo { };

in

pkgs.mkShell rec {
	buildInputs = with ourPkgs; [
		nodejs
		niv
		go
		gopls
		tinygo
		caddy
		nixos-generators
		jq
		qemu
	];

	shellHook = ''
		export PATH="$PATH:${PROJECT_ROOT}/node_modules/.bin"
		export NIX_PATH="''${NIX_PATH:-"$HOME/.nix-defexpr/channels"}:libdb.so=${PROJECT_ROOT}"
	'';

	PROJECT_ROOT = builtins.toString ./.;
}
