{ linuxKernel, linuxPackagesFor, writeText, runCommandLocal, pkgs, lib }:

let base = linuxKernel.kernels.linux_6_1; # https://endoflife.date/linux

	kernel = base.override {
		argsOverride.structuredExtraConfig = with lib.kernel; {
			LOCALVERSION = freeform "-v86-lto";

			M686 = yes;
			X86 = yes;
			X86_GENERIC = yes;
			X86_32 = yes;

			CONFIG_MODULE_SIG = no;
			CONFIG_MODULE_COMPRESS_NONE = yes;
			CONFIG_X86_MCE = no; # what machine?

			CONFIG_KERNEL_UNCOMPRESSED = yes;
			CONFIG_KERNEL_GZIP = no;
			CONFIG_INITRAMFS_COMPRESSION_NONE = yes;
		};
	};

in {
	inherit kernel;
	kernelPackage = linuxPackagesFor kernel;
}
