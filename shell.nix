{ pkgs ? import <nixpkgs> {} }:

let sources = import ./nix/sources.nix;
	mkShell = pkgs.mkShell;
in

let pkgs = import sources.nixpkgs {
		overlays = import ./nix/overlays.nix;
	};

	lib = pkgs.lib;

	tinygo = pkgs.callPackage ./nix/packages/tinygo { };

in

mkShell rec {
	buildInputs = with pkgs; [
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
