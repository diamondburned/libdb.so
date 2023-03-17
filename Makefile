.PHONY: all clean dev

SITE      = $(shell find site -type f) node_modules vite.config.ts package*.json
NIXOS     = $(shell find nixos -type f) $(shell find nix -type f)
TOOLCHAIN = $(shell find toolchain -type f)

# phony

all: dist

dev: nixos/result
	vite dev

clean:
	rm -r build

# real

dist: nixos/result $(SITE)
	vite build

nixos/result: nixos/packages/kernel/v86.base.config $(NIXOS) $(TOOLCHAIN)
	cd nixos && nix-build

nixos/packages/kernel/v86.base.config:
	cd nixos/packages/kernel && make
