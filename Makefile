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

dist: nixos/result $(SITE) $(NIXOS) $(TOOLCHAIN)
	vite build

nixos/result: $(NIXOS) $(TOOLCHAIN)
	cd nixos && nix-build
