{ pkgs, lib, system ? builtins.currentSystem }:

with lib;
with builtins;

let toybox = pkgs.toybox;
	kernel = pkgs.linuxKernel.kernels.linux_6_1;

	combineConfigs = paths: pkgs.runCommandLocal "v86-kernel.config" {} ''
		paths=( ${concatStringsSep " " paths} )
		for path in "''${paths[@]}"; do
			while IFS= read -r line; do
				if [[ "$line" == "CONFIG_"* ]]; then
					echo "$line" >> $out
				fi
			done < "$path"
		done
	'';

	kernelConfig = combineConfigs [
		./v86.config
		./v86.browser-vm.config
	];
in

pkgs.runCommand "toybox-exp" {} ''
	cp -r ${toybox.src} .
	cp -r ${kernel.src} linux/
	cp ${kernelConfig} linux/.config
''
