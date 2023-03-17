{ config, lib, pkgs, ... }:

{
	options.boot.initrd.services.swraid = {
		enable = lib.mkEnableOption "no-op";

		mdadmConf = lib.mkOption {
			description = lib.mdDoc "Contents of {file}`/etc/mdadm.conf` in initrd.";
			type = lib.types.lines;
			default = "";
		};
	};
}
