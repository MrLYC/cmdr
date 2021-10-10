ROOT_DIR ?= .
BIN_DIR ?= ${ROOT_DIR}/bin
TARGET ?= ${BIN_DIR}/cmdr

.PHONY: build
build:
	go build -o "${TARGET}" .