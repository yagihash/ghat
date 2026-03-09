#!/bin/bash
set -euo pipefail

# Only run in remote (web) environments
if [ "${CLAUDE_CODE_REMOTE:-}" != "true" ]; then
  exit 0
fi

# Read required Go version from go.mod (e.g. "1.26.1")
GO_VER=$(grep "^go " "${CLAUDE_PROJECT_DIR}/go.mod" | awk '{print $2}')
INSTALL_DIR="/usr/local/go${GO_VER}"
TOOLCHAIN_MODULE="golang.org/toolchain@v0.0.1-go${GO_VER}.linux-amd64"

# Download and install the required Go version if not already present
if [ ! -d "$INSTALL_DIR" ]; then
  echo "Installing Go ${GO_VER}..."
  ZIP="/tmp/go${GO_VER}.zip"
  EXTRACT_DIR="/tmp/go${GO_VER}_extract"

  curl -fsSL -o "$ZIP" \
    "https://proxy.golang.org/golang.org/toolchain/@v/v0.0.1-go${GO_VER}.linux-amd64.zip"

  mkdir -p "$EXTRACT_DIR"
  unzip -q "$ZIP" -d "$EXTRACT_DIR"
  mv "${EXTRACT_DIR}/${TOOLCHAIN_MODULE}" "$INSTALL_DIR"
  chmod +x "${INSTALL_DIR}/bin/go"
  rm -rf "$ZIP" "$EXTRACT_DIR"

  echo "Go ${GO_VER} installed at ${INSTALL_DIR}"
fi

# Point /usr/local/go at the required version
ln -sfn "$INSTALL_DIR" /usr/local/go

# Prevent go from trying to auto-download toolchains
echo 'export GOTOOLCHAIN=local' >> "$CLAUDE_ENV_FILE"

echo "Go toolchain ready: $(${INSTALL_DIR}/bin/go version)"
