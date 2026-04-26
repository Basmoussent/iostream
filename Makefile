BINARY      := iostream
GUI_BINARY  := iostream-gui
CLI_PKG     := ./cmd/$(BINARY)
GUI_PKG     := ./cmd/$(GUI_BINARY)
BIN_DIR     := bin
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS     := -s -w -X main.version=$(VERSION)
GUI_LDFLAGS := $(LDFLAGS) -H windowsgui

# Two binaries:
#   iostream      → the CLI (pure-Go default, cgo opt-in for libusb backend)
#   iostream-gui  → WebView wrapper, always cgo (uses WebView2/WebKit)
#
# Plain 'build' / 'vet' / 'windows' targets skip the GUI so they keep working
# on hosts without webview headers. The *-cgo targets and the explicit gui*
# targets pull it in.

.PHONY: all build cgo gui gui-windows-cgo gui-linux-cgo gui-darwin-cgo \
        windows windows-cgo linux linux-cgo darwin darwin-cgo \
        run test vet fmt tidy clean

all: build

build:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY) $(CLI_PKG)

cgo: build
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY) $(CLI_PKG)

gui:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(GUI_BINARY) $(GUI_PKG)

windows:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY).exe $(CLI_PKG)

windows-cgo:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY).exe $(CLI_PKG)

gui-windows-cgo:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags '$(GUI_LDFLAGS)' -o $(BIN_DIR)/$(GUI_BINARY).exe $(GUI_PKG)

linux:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY)-linux $(CLI_PKG)

linux-cgo:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY)-linux $(CLI_PKG)

gui-linux-cgo:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(GUI_BINARY)-linux $(GUI_PKG)

darwin:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY)-darwin $(CLI_PKG)

darwin-cgo:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY)-darwin $(CLI_PKG)

gui-darwin-cgo:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(GUI_BINARY)-darwin $(GUI_PKG)

run: build
	@$(BIN_DIR)/$(BINARY) $(ARGS)

test:
	go test ./cmd/iostream/... ./internal/...

vet:
	CGO_ENABLED=0 go vet ./cmd/iostream/... ./internal/...

fmt:
	gofmt -s -w .

tidy:
	go mod tidy

clean:
	rm -rf $(BIN_DIR)
