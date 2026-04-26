BINARY      := iphone-mirror
PKG         := ./cmd/$(BINARY)
BIN_DIR     := bin
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS     := -s -w -X main.version=$(VERSION)

# Native builds default to CGO_ENABLED=0 → the libusb backend stub is used.
# Real device I/O needs CGO_ENABLED=1 and libusb-1.0 dev headers; the cgo
# targets below opt in.

.PHONY: all build cgo windows windows-cgo linux linux-cgo darwin darwin-cgo \
        run test vet fmt tidy clean

all: build

build:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY) $(PKG)

cgo:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY) $(PKG)

windows:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY).exe $(PKG)

windows-cgo:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY).exe $(PKG)

linux:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY)-linux $(PKG)

linux-cgo:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY)-linux $(PKG)

darwin:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY)-darwin $(PKG)

darwin-cgo:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY)-darwin $(PKG)

run: build
	@$(BIN_DIR)/$(BINARY) $(ARGS)

test:
	go test ./...

vet:
	CGO_ENABLED=0 go vet ./...

fmt:
	gofmt -s -w .

tidy:
	go mod tidy

clean:
	rm -rf $(BIN_DIR)
