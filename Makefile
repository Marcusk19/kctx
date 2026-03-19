VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X github.com/Marcusk19/kctx/internal/cli.version=$(VERSION)"

.PHONY: build install clean test

build:
	go build $(LDFLAGS) -o kctx ./cmd/kctx

install:
	go install $(LDFLAGS) ./cmd/kctx

clean:
	rm -f kctx

test:
	go test ./...
