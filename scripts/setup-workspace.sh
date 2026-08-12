#!/usr/bin/env bash
# Generate go.work from go.work.example for local development.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORK_FILE="$SCRIPT_DIR/go.work"

if [ -f "$WORK_FILE" ]; then
	echo "go.work already exists, skipping generation"
	exit 0
fi

EXAMPLE_FILE="$SCRIPT_DIR/go.work.example"
if [ ! -f "$EXAMPLE_FILE" ]; then
	echo "ERROR: go.work.example not found at $EXAMPLE_FILE" >&2
	exit 1
fi

cp "$EXAMPLE_FILE" "$WORK_FILE"
echo "Generated go.work from go.work.example"
