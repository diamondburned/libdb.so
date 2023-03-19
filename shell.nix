{ pkgs ? import <nixpkgs> {} }:

let sources = import ./nix/sources.nix;
	ourPkgs = import sources.nixpkgs {
		overlays = import ./nix/overlays.nix;
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
		jq
		openssl
		gomod2nix
		# tinygo
	];

	shellHook = ''
		export PATH="$PATH:${PROJECT_ROOT}/node_modules/.bin"
		export NIX_PATH="''${NIX_PATH:-"$HOME/.nix-defexpr/channels"}:libdb.so=${PROJECT_ROOT}"
	'';

	PROJECT_ROOT = builtins.toString ./.;
}
