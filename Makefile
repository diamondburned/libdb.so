ROOT = $(shell pwd)
SITE = $(shell find site -type f) node_modules vite.config.ts package*.json
PUBLIC = $(shell find public/_fs -type f) $(shell find public/_assets -type f)
GOFILES = $(shell find vm -type f) go.mod go.sum
GOMODULE = $(shell go list -m)

# phony

.PHONY: all
all: build/dist

.PHONY: dev
dev: dist-deps
	vite dev

.PHONY: dist-deps
dist-deps: \
	$(SITE) \
	site/components/Terminal/color-schemes.json \
	build/public \
	build/vm.wasm

.PHONY: clean
clean:
	rm -r build

.PHONY: jsonld
jsonld: $(shell find jsonld -type f)
	tsx jsonld/_render.ts > public/_fs/0xd14.jsonld

.PHONY: openapi
openapi: backend/openapi

# real paths

OPENAPI_SCHEMAS = $(shell find backend/openapi -name '*.yml' -a -not -name '_*')
OPENAPI_CODE    = $(patsubst backend/openapi/%.yml,backend/openapi/%api,$(OPENAPI_SCHEMAS))

backend/openapi/%api: backend/openapi/%.yml build/openapi.config.yml
	mkdir -p "$@"
	cd backend/openapi && oapi-codegen \
		--config $(ROOT)/build/openapi.config.yml \
		--package "$(shell basename $@)" \
		-o "$(shell basename $@)/openapi.gen.go" "$(shell basename $<)"

backend/openapi: $(OPENAPI_CODE)

site/components/Terminal/color-schemes.json: ./scripts/xtermjs-colors
	./scripts/xtermjs-colors > $@

node_modules: package-lock.json package.json
	npm install

build/public/_fs: $(PUBLIC)
	rm -rf $@ && mkdir -p $@
	cp -r public/_fs/* $@

build/openapi.config.yml: backend/openapi/_config.yml $(OPENAPI_SCHEMAS)
	./scripts/import-mapping "$<" > "$@"

build/public/_fs.json: $(PUBLIC) build/public/_fs scripts/jsonfs
	cd $(dir $@) && bash $(ROOT)/scripts/jsonfs _fs > _fs.json

build/public: build/public/_fs build/public/_fs.json
	cp -r ./public/_assets $@

build/dist: dist-deps
	vite build

build/vm.wasm: $(GOFILES)
	mkdir -p build
	GOOS=js GOARCH=wasm go build -o build/vm.wasm ./vm/cmd/vm-wasm

build/backend: $(GOFILES) backend/openapi
	mkdir -p build
	go build -o build/backend ./backend/cmd/backend

#build/console.wasm: build/wasm_exec.js $(GOFILES)
#	mkdir -p build
#	tinygo build -o build/console.wasm -opt 2 -scheduler asyncify -target wasm ./cmd/console-wasm/

build/wasm_exec.js:
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js $@
