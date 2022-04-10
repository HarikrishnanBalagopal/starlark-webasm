# BIN_PATH=public/assets/wasm/main.wasm
BIN_PATH=main.wasm

.PHONY: build-prod
build-prod:
	GOOS=js GOARCH=wasm go build -ldflags "-s -w" -o "${BIN_PATH}"

.PHONY: build
build:
	GOOS=js GOARCH=wasm go build -o "${BIN_PATH}"

.PHONY: run
run:
	cd public && python3 -m http.server 8080

.PHONY: release-dry-run
release-dry-run:
	npm pack --dry-run
