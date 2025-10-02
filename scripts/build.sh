#!/bin/bash

export CGO_ENABLED=0

COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo unknown)

echo "Current commit: "$COMMIT

go build -ldflags="-X github.com/YardRat0117/foxbox/internal/version.Commit=$COMMIT" -o bin/foxbox ./cmd/foxbox
