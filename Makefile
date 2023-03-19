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

build/console.wasm: build/wasm_exec.js $(GOFILES)
	mkdir -p build
	GOOS=js GOARCH=wasm go build -o build/console.wasm ./cmd/console-wasm
# tinygo build -o build/console.wasm -opt 2 -scheduler asyncify -target wasm ./cmd/console-wasm/

build/wasm_exec.js:
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js $@
