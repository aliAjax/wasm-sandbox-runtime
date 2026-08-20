GO ?= /tmp/go-toolchain/go/bin/go
all: fmt test vet build
fmt:
	$(GO)fmt -w $$(find . -name '*.go')
test:
	$(GO) test ./...
race:
	$(GO) test -race ./...
vet:
	$(GO) vet ./...
build:
	$(GO) build ./...
run:
	$(GO) run ./cmd/sandboxd
smoke:
	./scripts/smoke.sh
