pkgs: {
	v86 = (pkgs.callPackage ./v86.nix {}).kernelPackage;
	uncompressed = (pkgs.callPackage ./uncompressed.nix {}).kernelPackage;
}
