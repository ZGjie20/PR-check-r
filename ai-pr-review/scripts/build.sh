#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

mkdir -p bin
go build -o bin/api ./cmd/api
go build -o bin/cli ./cmd/cli

echo "Built bin/api and bin/cli"
