with builtins;
with (import <nixpkgs> {}).lib;

let sources = import <libdb.so/nix/sources.nix>;
	nixpkgs = sources.nixpkgs;
	pkgs = import nixpkgs {};

	nixos' = import "${nixpkgs}/nixos" {
		configuration = import ./configuration.nix;
		system = "i686-linux"; # copyv86 architecture
	};

	nixos = with nixos'; trace ''Build Information:
  initrd kernel modules:
    required:  ${toString config.boot.initrd.kernelModules}
    available: ${toString config.boot.initrd.availableKernelModules}
  kernel modules:
    ${toString config.boot.kernelModules}
  packages:
    ${toString (unique (map (v: v.name) config.environment.systemPackages))}
'' nixos';

in pkgs.callPackage "${nixpkgs}/nixos/lib/make-disk-image.nix" rec {
	inherit (nixos) config;
	fsType = config.fileSystems."/".fsType;
	copyChannel = false;
	additionalSpace = "0M";
	partitionTableType = "legacy";
}
