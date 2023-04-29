ROOT_DIR ?= .
BIN_DIR ?= ${ROOT_DIR}/bin
TARGET ?= ${BIN_DIR}/cmdr

${TARGET}:
	go build -o "${TARGET}" .

.PHONY: build
build: ${TARGET}

.PHONY: goreleaser
goreleaser:
	goreleaser build --skip-validate --single-target --rm-dist --snapshot

.PHONY: test
test:
	go test -gcflags=all=-l ./...

.PHONY: generate
generate:
	go generate ./...

usage-test.sh: ${TARGET}
	bash usage-test.sh