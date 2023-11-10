build:
	rm -rf bin
	go build -o bin/terra ./cmd/terra
	strip -s bin/terra

build-musl:
	CGO_ENABLED=1 CC=musl-gcc go build \
		--ldflags 'linkmode external -extldflags="-static"' -o bin/terra ./cmd/terra
	strip -s bin/terra

run:
	go run ./cmd/terra

install:
	make build
	sudo cp -v bin/terra /usr/local/bin/terra
