let sources = import ./sources.nix;
in

[
	(self: super: {
		go = super.go_1_20;
		buildGoModule = super.buildGo120Module;
	})
	(self: super: {
		npmlock2nix = super.callPackage sources.npmlock2nix { };
	})
	(import "${sources.gomod2nix}/overlay.nix")
]
