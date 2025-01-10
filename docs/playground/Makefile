.PHONY: build
build: build/wasm build/page

.PHONY: build/wasm
build/wasm:
	GOOS=js GOARCH=wasm go build -o src/yaml.wasm ./cmd/yaml

.PHONY: build/page
build/page:
	npm run build

.PHONY: deps
deps:
	npm ci

.PHONY: dev
dev:
	npm run dev
