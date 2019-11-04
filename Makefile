.PHONY: wasm

wasm:
	GOOS=js GOARCH=wasm go build -o docs/yaml.wasm docs/wasm.go
