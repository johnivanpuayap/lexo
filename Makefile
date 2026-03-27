.PHONY: build test clean wasm wasm-exec

# Build WASM binary
wasm:
	GOOS=js GOARCH=wasm go build -o web/public/lexo.wasm ./cmd/wasm

# Copy Go's WASM JS runtime helper
wasm-exec:
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/public/wasm_exec.js

# Run all Go tests
test:
	go test ./pkg/... -v

# Build everything
build: wasm-exec wasm

# Clean build artifacts
clean:
	rm -f web/public/lexo.wasm web/public/wasm_exec.js
