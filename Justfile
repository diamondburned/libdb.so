build: build-spiral

build-spiral:
	tinygo build -o dist/spiral.wasm -target wasm ./vm/cmd/spiral-wasm/
