SCRIPTS_DIR ?= $(HOME)/Development/github.com/rios0rios0/pipelines
-include $(SCRIPTS_DIR)/makefiles/common.mk
-include $(SCRIPTS_DIR)/makefiles/golang.mk

.PHONY: build debug build-musl run install

build:
	mkdir -p bin && rm -rf bin/terra
	go mod tidy
	go build -o bin/terra ./cmd/terra
	strip -s bin/terra

debug:
	rm -rf bin
	go build -gcflags "-N -l" -o bin/terra ./cmd/terra

build-musl:
	CGO_ENABLED=1 CC=musl-gcc go build \
		--ldflags 'linkmode external -extldflags="-static"' -o bin/terra ./cmd/terra
	strip -s bin/terra

run:
	go run ./cmd/terra

install:
	make build
	mkdir -p ~/.local/bin
	cp -v bin/terra ~/.local/bin/terra
