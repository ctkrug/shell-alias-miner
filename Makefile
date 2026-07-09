.PHONY: build test test-js vet fmt site wasm clean

# Everything except cmd/wasm builds and tests under the host GOOS/GOARCH.
build:
	go build ./internal/...

test:
	go test ./internal/...

# Pure DOM-free logic in site/main.js (thresholds, the explain formula).
test-js:
	node --test site/main.test.js

vet:
	go vet ./internal/...
	GOOS=js GOARCH=wasm go vet ./cmd/wasm

# Fails (non-zero exit, prints offending files) if anything is unformatted.
fmt:
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt needed on:"; echo "$$unformatted"; exit 1; \
	fi

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
