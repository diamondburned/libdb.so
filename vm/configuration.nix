{ config, pkgs, lib, modulesPath, ... }:

{
	imports = [
		(modulesPath + "/profiles/minimal.nix")
	];

	# TODO: build a custom minimal kernel
}
