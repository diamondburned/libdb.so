let sources = import ./nix/sources.nix;
	ourPkgs = import sources.nixpkgs {
		overlays = import ./nix/overlays.nix;
	};

	systemPkgs = with builtins.tryEval <nixpkgs>;
		if success then
			import value {}
		else
			ourPkgs;
in

{ pkgs ? systemPkgs }:

let sources = import ./nix/sources.nix;
	mkShell = pkgs.mkShell;
in

# prefer our Nixpkgs
let pkgs = ourPkgs;
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
