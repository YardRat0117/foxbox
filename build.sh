#!/bin/bash

COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo unknown)

echo "Current commit: "$COMMIT

go build -ldflags="-X github.com/YardRat0117/foxbox/src/version.Commit=$COMMIT" -o foxbox ./src

