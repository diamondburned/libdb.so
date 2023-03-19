let sources = import ./nix/sources.nix;
in

{
	pkgs ? import sources.nixpkgs {},
	lib ? pkgs.lib,
	src ? builtins.filterSource
		(path: type:
			(baseNameOf path != ".git") &&
			(baseNameOf path != "build"))
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
		let module = pkgs.buildGoModule {
			pname = "libdb.so-vm-wasm";
			inherit version src;
	
			subPackages = [ "cmd/vm-wasm" ];
			vendorSha256 = "sha256:1mwymf0x569frm906xldn0qlik41mfif7w0avhwmzmd07npg0xnd";
	
			doCheck = false; # none to run
			CGO_ENABLED = 0;

			postInstall = ''
				mv $out/bin/js_wasm/vm-wasm $out/bin/vm.wasm
				rmdir $out/bin/js_wasm
			'';
		};
		in module.overrideAttrs (old: old // {
			GOOS = "js";
			GOARCH = "wasm";
		});

	npmPackage = pkgs.buildNpmPackage {
		pname = "libdb.so-vm-wasm";
		inherit version src;

		npmBuildScript = "build";
		npmDepsHash = "sha256:0pdpg76kccgj8dhf203mraxz9ys9aarjp41v3qj9p4mhqyqy3nbg";
		dontNpmInstall = true;

		preBuild = ''
			mkdir -p build
			cp ${vmWasm}/bin/vm.wasm build/vm.wasm
		'';

		installPhase = ''
			cp -r build/dist $out
		'';
	};
in

npmPackage
