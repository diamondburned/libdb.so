{ pkgs ? import <nixpkgs> {} }:

let sources = import ./nix/sources.nix;
	ourPkgs = import sources.nixpkgs { };

	lib = pkgs.lib;

	tinygo = ourPkgs.callPackage ./nix/packages/tinygo { };

in

pkgs.mkShell rec {
	buildInputs = with ourPkgs; [
		nodejs
		niv
		# We're stuck with Go 1.19 until we get tinygo 1.27.0, which I cannot
		# for the life of me get to build.
		go_1_19
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
