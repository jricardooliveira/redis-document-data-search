#!/bin/bash

# Validate that all required tools for the project are installed and available in PATH.
# Exits with non-zero status and a clear error message if any are missing.

set -e

declare -A TOOL_EXPLANATIONS
TOOL_EXPLANATIONS=(
  ["go"]="Go: Required to build and run the backend and CLI applications."
  ["redis-cli"]="redis-cli: Used for interacting with Redis/Valkey directly, especially for listing keys."
  ["jq"]="jq: Used for parsing and extracting fields from JSON in shell scripts."
  ["curl"]="curl: Used for making HTTP requests to the API endpoints from scripts."
  ["make"]="make: Used to build, test, and manage the project via the Makefile."
  ["shuf"]="shuf: Used in scripts to randomly sample keys or lines from files."
)

REQUIRED_TOOLS=("go" "redis-cli" "jq" "curl" "make" "shuf")
MISSING=()

for TOOL in "${REQUIRED_TOOLS[@]}"; do
  EXPLANATION="${TOOL_EXPLANATIONS[$TOOL]}"
  if command -v "$TOOL" >/dev/null 2>&1; then
    echo "[OK] $TOOL: $EXPLANATION"
  else
    echo "[MISSING] $TOOL: $EXPLANATION"
    MISSING+=("$TOOL")
  fi
done

if [ ${#MISSING[@]} -ne 0 ]; then
  echo -e "\n[ERROR] The following required tools are missing from your PATH:"
  for TOOL in "${MISSING[@]}"; do
    echo "  - $TOOL"
  done
  echo -e "\nPlease install the missing tools and try again."
  exit 1
else
  echo "\nAll required tools are installed."
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
