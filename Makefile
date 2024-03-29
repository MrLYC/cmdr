ROOT_DIR ?= .
BIN_DIR ?= ${ROOT_DIR}/bin
TARGET ?= ${BIN_DIR}/cmdr

.PHONY: build
build:
	go build -o "${TARGET}" .

.PHONY: goreleaser
goreleaser:
	goreleaser build --skip-validate --single-target --clean

.PHONY: test
test:
	go test -gcflags=all=-l ./...

.PHONY: generate
generate:
	go generate ./...
