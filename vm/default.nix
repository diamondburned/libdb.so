let sources = import <libdb.so/nix/sources.nix>;
	nixpkgs = sources.nixpkgs;
	pkgs = import nixpkgs {};

in pkgs.callPackage "${nixpkgs}/nixos/lib/make-disk-image.nix" {
	config = import "${nixpkgs}/nixos" {
		configuration = import ./configuration.nix;
		system = "i686-linux"; # copyv86 architecture
	};
	bootSize = "64M";
	copyChannel = false;
	additionalSpace = "0M";
	partitionTableType = "efi";
}
