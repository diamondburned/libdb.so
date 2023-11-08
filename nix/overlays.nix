let sources = import ./sources.nix;
in

[
	(self: super: {
		go = super.go_1_21;
		buildGoModule = super.buildGo121Module;
	})
	(self: super: {
		npmlock2nix = super.callPackage sources.npmlock2nix { };
	})
	(import "${sources.gomod2nix}/overlay.nix")
]
