# ==== Config ====
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
COMMIT ?= $(shell git rev-parse --short HEAD)
PKG_VERSION ?= main
MAIN := ./cmd/foxbox
RELEASE_BIN ?= bin/foxbox-$(GOOS)-$(GOARCH)

# ==== Build Info ====
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)

# ==== Targets ====

all: build compress

# Compile Go binary
build:
	@echo "Building $(RELEASE_BIN) for $(GOOS)/$(GOARCH)..."
	@echo "Commit: $(COMMIT)"
	@mkdir -p $(dir $(RELEASE_BIN))
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags="-s -w -X $(PKG_VERSION).Commit=$(COMMIT)" \
		-o $(RELEASE_BIN) $(MAIN)
	@du -h $(RELEASE_BIN)

# Compress with UPX
compress:
	@echo "Compressing with UPX..."
	upx --best --lzma $(RELEASE_BIN) || echo "UPX not found or compression failed."
	@du -h $(RELEASE_BIN)

# Cleaning up
clean:
	@echo "Cleaning..."
	rm -rf bin

.PHONY: all build compress clean

