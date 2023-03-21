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

let mkShell = pkgs.mkShell;
in

# prefer our Nixpkgs
let pkgs = ourPkgs;
	lib = pkgs.lib;

	tinygo = pkgs.callPackage ./nix/packages/tinygo { };

	# We always want gopls to run in GOOS=js mode so we can get proper autocompletions.
	gopls = pkgs.writeShellScriptBin "gopls" ''
		GOOS=js GOARCH=wasm exec ${pkgs.gopls}/bin/gopls "$@"
	'';

in

mkShell rec {
	buildInputs = with pkgs; [
		nodejs
		niv
		go
		gopls
		jq
		gomod2nix
		# tinygo
	];

	# GOOS = "js";
	# GOARCH = "wasm";

	shellHook = ''
		export PATH="$PATH:${PROJECT_ROOT}/node_modules/.bin"
		export NIX_PATH="''${NIX_PATH:-"nixpkgs=${sources.nixpkgs}"}:libdb.so=${PROJECT_ROOT}"
	'';

	PROJECT_ROOT = builtins.toString ./.;
}
