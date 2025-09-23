SCRIPTS_DIR := $(HOME)/Development/github.com/rios0rios0/pipelines
REPO_URL    := https://github.com/rios0rios0/pipelines.git

.PHONY: all scripts lint lint-fix horusec test

all: lint horusec test

scripts:
	if [ ! -d "$(SCRIPTS_DIR)" ]; then \
	  git clone $(REPO_URL) $(SCRIPTS_DIR); \
	else \
	  cd $(SCRIPTS_DIR) && git pull; \
	fi

lint: scripts
	$(SCRIPTS_DIR)/global/scripts/golangci-lint/run.sh .

lint-fix: scripts
	$(SCRIPTS_DIR)/global/scripts/golangci-lint/run.sh --fix .

horusec: scripts
	$(SCRIPTS_DIR)/global/scripts/horusec/run.sh .

test: scripts
	$(SCRIPTS_DIR)/global/scripts/golang/test/run.sh .

clear:
	sudo rm -rf .go
	rm -rf bin coverage.* cobertura.xml junit.xml

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
	cp -v bin/terra ~/.local/bin/terra
