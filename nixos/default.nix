with builtins;
with (import <nixpkgs> {}).lib;

let sources = import <libdb.so/nix/sources.nix>;
	nixpkgs = sources.nixpkgs;
	pkgs = import nixpkgs {};

	evalConfig = import <nixpkgs/nixos/lib/eval-config.nix>;

	config' = evalConfig {
		modules = [ (import ./configuration.nix) ];
		system = "i686-linux"; # copyv86 architecture
	};

	# nixos' = import "${nixpkgs}/nixos" {
	# 	configuration = import ./configuration.nix;
	# 	system = "i686-linux"; # copyv86 architecture
	# };

in config'

# 	config = with config'; trace ''Build Information:
#   initrd kernel modules:
#     required:  ${toString boot.initrd.kernelModules}
#     available: ${toString boot.initrd.availableKernelModules}
#   kernel modules:
#     ${toString boot.kernelModules}
#   packages:
#     ${toString (unique (map (v: v.name) environment.systemPackages))}
# '' config';

# in import "${nixpkgs}/nixos/lib/make-disk-image.nix" rec {
# 	inherit config pkgs;
# 	inherit (pkgs) lib;
#
# 	fsType = config.fileSystems."/".fsType;
# 	format = "raw";
# 	diskSize = "auto";
# 	additionalSpace = "0M";
# 	installBootLoader = true;
# 	partitionTableType = "legacy";
# 	copyChannel = false;
# 	memSize = 1024;
# }
