{ linuxKernel, writeText, runCommandLocal, lib }:

with lib;
with builtins;

let base = linuxKernel.kernels.linux_6_1; # https://endoflife.date/linux

	baseConfig = linuxKernel.linuxConfig {
		makeTarget = "tinyconfig";
		src = base.src;
	};

	combineConfigs = paths: runCommandLocal "v86-kernel.config" {} ''
		paths=( ${concatStringsSep " " paths} )
		for path in "''${paths[@]}"; do
			while IFS= read -r line; do
				if [[ "$line" == "CONFIG_"* ]]; then
					echo "$line" >> $out
				fi
			done < "$path"
		done
	'';

	kernelPackage = linuxKernel.customPackage {
		inherit (base) src version;
		configfile = combineConfigs [
			./v86.config
			./v86.base.config
		];
		allowImportFromDerivation = true;
	};

in {
	inherit kernelPackage baseConfig;
}
