let sources = import ./nix/sources.nix;
in

{
	pkgs ? import sources.nixpkgs {
		overlays = import ./nix/overlays.nix;
	},
	lib ? pkgs.lib,
	src ? builtins.filterSource
		(path: type:
			(baseNameOf path != ".git") &&
			(baseNameOf path != "build") &&
			(baseNameOf path != "node_modules"))
		(./.),
	version ? "git",
	outputHash ? "sha256:${lib.fakeSha256 }",
}:

let stdenv = pkgs.stdenv;

	buildGoWasmModule = pkgs.buildGoModule.override {
		go = pkgs.go // {
			GOOS = "js";
			GOARCH = "wasm";
		};
	};

	vmWasm =
		let module = pkgs.buildGoApplication {
			inherit version src;
			pname = "libdb.so-vm-wasm";
			go = pkgs.go_1_20;
			modules = ./gomod2nix.toml;
			subPackages = [ "cmd/vm-wasm" ];

			CGO_ENABLED = 0;
			doCheck = false; # none to run

			postInstall = ''
				mv $out/bin/js_wasm/vm-wasm $out/bin/vm.wasm
				rmdir $out/bin/js_wasm
			'';
		};
		in module.overrideAttrs (old: old // {
			GOOS = "js";
			GOARCH = "wasm";
		});

	nodeModules = pkgs.npmlock2nix.v2.node_modules {
		inherit src;
		# mkDerivation hates us because we have a Makefile. We'll override
		# installPhase to fix that.
		installPhase = "mv node_modules $out/";
	};
in

stdenv.mkDerivation {
	inherit version src;
	pname = "libdb.so";

	buildInputs = with pkgs; [
		coreutils
		bash
		jq
		# go
		nodejs
	];

	preBuild = ''
		set -x

		mkdir -p build
		cp -r ${vmWasm}/bin/vm.wasm build/vm.wasm

		cp -r ${nodeModules} node_modules
		chown -R $(id -u):$(id -g) node_modules
		chmod -R +w node_modules
		export PATH="$PATH:$PWD/node_modules/.bin"

		set +x
	'';

	installPhase = ''
		cp -r build/dist $out
	'';
}
