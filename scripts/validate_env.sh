#!/bin/bash

# Validate that all required tools for the project are installed and available in PATH.
# Exits with non-zero status and a clear error message if any are missing.

set -e

REQUIRED_TOOLS=(
  "go"
  "redis-cli"
  "jq"
  "curl"
  "make"
  "shuf"
)

MISSING=()

for TOOL in "${REQUIRED_TOOLS[@]}"; do
  if ! command -v "$TOOL" >/dev/null 2>&1; then
    MISSING+=("$TOOL")
  fi
done

if [ ${#MISSING[@]} -ne 0 ]; then
  echo "\n[ERROR] The following required tools are missing from your PATH:"
  for TOOL in "${MISSING[@]}"; do
    echo "  - $TOOL"
  done
  echo "\nPlease install the missing tools and try again."
  exit 1
else
  echo "All required tools are installed."
fi

# Optionally, check Go version
if command -v go >/dev/null 2>&1; then
  GOVERSION=$(go version | awk '{print $3}')
  GOVERSION_NUM=$(echo $GOVERSION | sed 's/go//')
  REQUIRED_GO="1.18"
  if [[ $(printf '%s\n' "$REQUIRED_GO" "$GOVERSION_NUM" | sort -V | head -n1) != "$REQUIRED_GO" ]]; then
    echo "[ERROR] Go version $REQUIRED_GO or higher is required. Found $GOVERSION_NUM."
    exit 1
  fi
fi
