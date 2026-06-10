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

.PHONY: coverage
coverage:
	go test -gcflags=all=-l ./... -coverprofile=coverage.out -covermode=count
	awk 'NR == 1 || ($$1 !~ /\/mock\// && $$1 !~ /_string\.go:/ && $$1 !~ /\/cmd\/internal\/testutils\// && $$1 !~ /\/main\.go:/)' coverage.out > coverage.business.out
	go tool cover -func=coverage.business.out

.PHONY: coverage-check
coverage-check: coverage
	@coverage=$$(go tool cover -func=coverage.business.out | awk '/^total:/ { sub(/%/, "", $$3); print $$3 }'); \
	awk -v coverage="$$coverage" 'BEGIN { if (coverage + 0 < 95.0) { printf "business coverage %.1f%% is below 95.0%%\n", coverage; exit 1 }; printf "business coverage %.1f%% meets 95.0%%\n", coverage }'

.PHONY: generate
generate:
	go generate ./...
