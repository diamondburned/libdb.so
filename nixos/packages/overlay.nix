self: super:

with super.lib;
with builtins;

{
	diamondburned = rec {
		kernelPackages = import ./kernel self;
	};
}
