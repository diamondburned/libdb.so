.PHONY: all clean dev

SITE    = $(shell find site -type f) node_modules vite.config.ts package*.json
GOFILES = $(shell find console -type f) $(shell find cmd -type f) go.mod go.sum

# phony

all: build/dist

dev: build/console.wasm
	vite dev

clean:
	rm -r build

# real

build/dist: build/console.wasm $(SITE)
	vite build

build/console.wasm: $(GOFILES)
	mkdir -p build
	tinygo build -o build/console.wasm -target wasm ./cmd/console-wasm/
