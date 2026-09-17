#!/bin/sh

# set -eu

# ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
# cd "$ROOT_DIR"

# go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -d ./cmd/server -g main.go -o ./docs
exec go run ./cmd/server