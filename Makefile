# ==== Config ====
BINARY := bin/foxbox
MAIN := ./cmd/foxbox
PKG_VERSION := github.com/YardRat0117/foxbox/internal/version

# ==== Build Info ====
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)

# ==== Targets ====

all: build compress

# 编译 Go 二进制
build:
	@echo "Building $(BINARY)..."
	@echo "Commit: $(COMMIT)"
	@mkdir -p bin
	go build -ldflags="-s -w -X $(PKG_VERSION).Commit=$(COMMIT)" -o $(BINARY) $(MAIN)
	@du -h $(BINARY)

# 压缩二进制
compress:
	@echo "Compressing with UPX..."
	upx --best --lzma $(BINARY) || echo "UPX not found or compression failed."
	@du -h $(BINARY)

# 运行
run: all
	@echo "Running $(BINARY)..."
	@./$(BINARY)

# 清理
clean:
	@echo "Cleaning..."
	rm -f $(BINARY)
	rm -rf bin

.PHONY: all build compress run clean

