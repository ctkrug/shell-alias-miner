.PHONY: build test vet wasm site clean

# Everything except cmd/wasm builds and tests under the host GOOS/GOARCH.
build:
	go build ./internal/...

test:
	go test ./...

vet:
	go vet ./internal/... ./cmd/...

# cmd/wasm only compiles under GOOS=js GOARCH=wasm (it imports syscall/js),
# so it's built separately from the rest of the module.
wasm:
	GOOS=js GOARCH=wasm go build -o site/main.wasm ./cmd/wasm

# wasm_exec.js is the JS glue Go ships for running its wasm output; its
# location has moved between Go versions, so try both.
site: wasm
	@cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" site/wasm_exec.js 2>/dev/null || \
	 cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" site/wasm_exec.js

clean:
	rm -f site/main.wasm site/wasm_exec.js
