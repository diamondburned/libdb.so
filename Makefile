.PHONY: all clean build-init

all: build

clean:
	rm -r build

build: build-init build/vm build/toolchain build/site

build-init:
	mkdir -p build

build/vm:
	# nix-build -o ./build/vm ./vm

build/toolchain:
	go build -v ./toolchain/...

build/site: site node_modules vite.config.ts package.json package-lock.json
	vite build --outDir $$PROJECT_ROOT/build/site
