build:
	mkdir -p bin && rm -rf bin/terra
	go install github.com/google/wire/cmd/wire@latest
	wire ./...
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
