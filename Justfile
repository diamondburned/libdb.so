goflags := ""

dev: build-spiral
	pnpm dev

build: build-spiral build-astro

build-astro:
	pnpm build

build-spiral $GOOS="js" $GOARCH="wasm":
	go build -o dist/spiral.wasm {{goflags}} -ldflags="-s -w" ./vm/cmd/spiral-wasm/
