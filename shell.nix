{ pkgs ? import <nixpkgs> {} }:

let sources = import ./nix/sources.nix;
	ourPkgs = import sources.nixpkgs { };

in

pkgs.mkShell rec {
	buildInputs = with ourPkgs; [
		nodejs
		niv
		go_1_20
		caddy
		nixos-generators
		jq
		qemu
	];

	shellHook = ''
		export PATH="$PATH:${PROJECT_ROOT}/node_modules/.bin"
		export NIX_PATH="''${NIX_PATH:-"$HOME/.nix-defexpr/channels"}:libdb.so=${PROJECT_ROOT}:libdb.so/nixpkgs=${sources.nixpkgs}"
	'';

	PROJECT_ROOT = builtins.toString ./.;
}
