.PHONY: all clean dev

SITE    = $(shell find site -type f) node_modules vite.config.ts package*.json
PUBLIC  = $(shell find public/_fs -type f)
GOFILES = $(shell find cmd internal vm -type f) go.mod go.sum

# phony

all: build/dist

dev: dist-deps
	vite dev

dist-deps: $(SITE) site/components/Terminal/color-schemes.json public/_fs.json build/vm.wasm

clean:
	rm -r build

# real

public/_fs.json: $(PUBLIC) ./scripts/jsonfs
	cd public && bash ../scripts/jsonfs _fs > _fs.json

build/dist: dist-deps
	vite build

site/components/Terminal/color-schemes.json: ./scripts/xtermjs-colors
	./scripts/xtermjs-colors > $@

build/vm.wasm: $(GOFILES)
	mkdir -p build
	GOOS=js GOARCH=wasm go build -o build/vm.wasm ./cmd/vm-wasm

#build/console.wasm: build/wasm_exec.js $(GOFILES)
#	mkdir -p build
#	tinygo build -o build/console.wasm -opt 2 -scheduler asyncify -target wasm ./cmd/console-wasm/

build/wasm_exec.js:
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js $@
