{ config, pkgs, lib, modulesPath, ... }:

let sources = import <libdb.so/nix/sources.nix>;
in

{
	imports = [
		(modulesPath + "/profiles/minimal.nix")
		(import ./modules)
	];

	nixpkgs.overlays = [
		(import ./packages/overlay.nix)
	];

	fileSystems."/" = {
		fsType = "ext4"; # required by make-disk-image
		device = "/dev/disk/by-label/nixos";
		noCheck = true;
		autoResize = false; # slow
	};

	# v86 uses HDA.
	boot.loader.grub.device = "/dev/hda";

	# These don't work with systemdMinimal. They work with the current systemd,
	# but we're disabling them anyway to be lightweight.
	systemd.oomd.enable = false;
	systemd.coredump.enable = false;
	services.timesyncd.enable = false; # who needs a clock anyway

	# Use our minimal v86 kernel.
	boot.kernelPackages = pkgs.diamondburned.kernelPackages.v86;

	# These fail on our v86 kernel, so we disable it.
	boot.initrd.systemd.enable = false;
	boot.initrd.services.swraid.enable = false;

	# Unnecessary.
	boot.enableContainers = false;
	boot.initrd.includeDefaultModules = false;
	environment.defaultPackages = lib.mkForce [ ];

	# Emulated environment, so no need for udev.
	boot.hardwareScan = false;
	# services.udev.enable = false;

	# These are taken from vikanezrimaya/nixos-super-minimal.
	services.udisks2.enable = false;
	services.nscd.enable = false;
	system.nssModules = lib.mkForce [];
}
