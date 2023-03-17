{ config, pkgs, lib, ... }:

{
	imports = [
		./no-swraid.nix
		./no-ext.nix
	];

	disabledModules = [
		"tasks/swraid.nix"
		"tasks/filesystems/ext.nix"
	];
}
