.PHONY: all clean dev

SITE    = $(shell find site -type f) node_modules vite.config.ts package*.json
GOFILES = $(shell find . -type f -name '*.go') go.mod go.sum

# phony

all: build/dist

dev: site/components/Terminal/color-schemes.json build/vm.wasm $(SITE)
	vite dev

clean:
	rm -r build

# real

build/dist: site/components/Terminal/color-schemes.json build/vm.wasm $(SITE)
	npm i
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
