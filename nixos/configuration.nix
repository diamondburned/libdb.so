{ config, pkgs, lib, modulesPath, ... }:

let sources = import <libdb.so/nix/sources.nix>;
in

{
	imports = [
		(modulesPath + "/profiles/minimal.nix")
		# (import ./modules)
	];

	nixpkgs.overlays = [
		(import ./packages/overlay.nix)
	];

	# fileSystems."/" = {
	# 	fsType = "ext4"; # required by make-disk-image
	# 	device = "/dev/hda";
	# 	autoFormat = true;
	# 	# noCheck = false;
	# 	# autoResize = true; # slow
	# };

	# We only need ext4 for the root filesystem.
	# boot.initrd.supportedFilesystems = lib.mkForce [ "ext4" ];

	# v86 uses HDA.
	# boot.loader.grub.device = "/dev/hda";

	# These don't work with systemdMinimal. They work with the current systemd,
	# but we're disabling them anyway to be lightweight.
	# systemd.oomd.enable = false;
	# systemd.coredump.enable = false;
	# services.timesyncd.enable = false; # who needs a clock anyway

	# Use our minimal v86 kernel.
	# boot.kernelPackages = pkgs.diamondburned.kernelPackages.v86;

	# We know exactly what modules our kernel will provide.
	# boot.initrd.availableKernelModules = lib.mkForce [ "autofs4" ];
	# boot.initrd.services.swraid.enable = false;

	# Our kernel doesn't boot :(
	# boot.kernelPackages = pkgs.diamondburned.kernelPackages.uncompressed;

	# boot.initrd.systemd.enable = true;

	# Unnecessary.
	boot.enableContainers = false;
	# boot.initrd.includeDefaultModules = false;
	# environment.defaultPackages = lib.mkForce [ ];

	# Emulated environment, so no need for udev.
	# boot.hardwareScan = false;
	# services.udev.enable = false;

	# These are taken from vikanezrimaya/nixos-super-minimal.
	# services.udisks2.enable = false;
	# services.nscd.enable = false;
	# system.nssModules = lib.mkForce [];

	users.users.root.password = "root";
	users.mutableUsers = false;

	# Disable sound.
	sound.enable = false;
	hardware.pulseaudio.enable = false;
}
