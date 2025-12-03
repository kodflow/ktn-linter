#!/bin/bash

# KTN-Linter Universal Installer
# Installs ktn-linter from GitHub releases and configures golangci-lint integration
# Can be used on any Go project

set -e

REPO="kodflow/ktn-linter"
BINARY_NAME="ktn-linter"
VERSION="${KTN_VERSION:-latest}"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   KTN-Linter Universal Installer      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

case "$OS" in
    linux|darwin)
        ;;
    *)
        echo -e "${RED}❌ Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

echo -e "${YELLOW}📦 Platform: ${OS}-${ARCH}${NC}"

# Determine installation directory
if [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
elif [ -n "$HOME" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"

    # Add to PATH if not already there
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        echo -e "${YELLOW}💡 Adding $INSTALL_DIR to PATH${NC}"
        echo ""
        echo -e "${YELLOW}Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):${NC}"
        echo -e "${GREEN}export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}"
        echo ""
    fi
else
    echo -e "${RED}❌ Cannot determine installation directory${NC}"
    exit 1
fi

echo -e "${YELLOW}📂 Installation directory: $INSTALL_DIR${NC}"

# Get latest release or specific version
echo -e "${YELLOW}🔍 Fetching release information...${NC}"

if [ "$VERSION" = "latest" ]; then
    LATEST_RELEASE=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest")
    VERSION=$(echo "$LATEST_RELEASE" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    DOWNLOAD_URL=$(echo "$LATEST_RELEASE" | grep "browser_download_url.*${OS}-${ARCH}" | cut -d '"' -f 4)
else
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}-${OS}-${ARCH}"
fi

if [ -z "$DOWNLOAD_URL" ]; then
    echo -e "${RED}❌ No release found for ${OS}-${ARCH}${NC}"
    echo -e "${YELLOW}💡 Building from source...${NC}"

    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go is not installed. Cannot build from source.${NC}"
        exit 1
    fi

    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    git clone "https://github.com/${REPO}.git"
    cd ktn-linter
    go build -o "$INSTALL_DIR/$BINARY_NAME" ./cmd/ktn-linter
    cd /
    rm -rf "$TEMP_DIR"

    echo -e "${GREEN}✅ Built and installed from source${NC}"
else
    echo -e "${GREEN}✅ Found version: ${VERSION}${NC}"

    # Download binary
    TEMP_FILE=$(mktemp)
    echo -e "${YELLOW}⬇️  Downloading...${NC}"

    if ! curl -L -o "$TEMP_FILE" "$DOWNLOAD_URL"; then
        echo -e "${RED}❌ Download failed${NC}"
        rm -f "$TEMP_FILE"
        exit 1
    fi

    # Install binary
    chmod +x "$TEMP_FILE"
    mv "$TEMP_FILE" "$INSTALL_DIR/$BINARY_NAME"

    echo -e "${GREEN}✅ Installed to $INSTALL_DIR/$BINARY_NAME${NC}"
fi

# Verify installation
if command -v "$BINARY_NAME" &> /dev/null; then
    VERSION_OUTPUT=$("$BINARY_NAME" --version 2>&1 || echo "version check skipped")
    echo -e "${GREEN}✅ Installation verified: ${VERSION_OUTPUT}${NC}"
else
    echo -e "${YELLOW}⚠️  Binary installed but not in PATH${NC}"
    echo -e "${YELLOW}   Run: export PATH=\"$INSTALL_DIR:\$PATH\"${NC}"
fi

# Configure golangci-lint integration
echo ""
echo -e "${BLUE}═══════════════════════════════════════${NC}"
echo -e "${BLUE}  golangci-lint Integration (Optional)${NC}"
echo -e "${BLUE}═══════════════════════════════════════${NC}"
echo ""

read -p "$(echo -e ${YELLOW}Configure golangci-lint integration in current project? [y/N]: ${NC})" -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    GOLANGCI_CONFIG=".golangci.yml"

    if [ ! -f "$GOLANGCI_CONFIG" ]; then
        echo -e "${YELLOW}Creating $GOLANGCI_CONFIG...${NC}"
        cat > "$GOLANGCI_CONFIG" <<'EOF'
run:
  timeout: 5m
  tests: false

linters:
  enable:
    - gofmt
    - govet
    - staticcheck
    - unused
    - gosimple
    - ineffassign

linters-settings:
  custom:
    ktn:
      path: ktn-linter
      description: KTN-Linter - Strict Go best practices
      original-url: https://github.com/kodflow/ktn-linter

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - ktn
EOF
        echo -e "${GREEN}✅ Created $GOLANGCI_CONFIG with ktn-linter${NC}"
    else
        echo -e "${YELLOW}⚠️  $GOLANGCI_CONFIG already exists${NC}"
        echo -e "${YELLOW}💡 Add this to your linters-settings.custom section:${NC}"
        echo ""
        echo -e "${GREEN}linters-settings:${NC}"
        echo -e "${GREEN}  custom:${NC}"
        echo -e "${GREEN}    ktn:${NC}"
        echo -e "${GREEN}      path: ktn-linter${NC}"
        echo -e "${GREEN}      description: KTN-Linter - Strict Go best practices${NC}"
        echo -e "${GREEN}      original-url: https://github.com/kodflow/ktn-linter${NC}"
        echo ""
    fi

    # Create Makefile targets if Makefile doesn't exist
    if [ ! -f "Makefile" ]; then
        echo -e "${YELLOW}💡 Creating Makefile with ktn-linter targets...${NC}"
        cat > "Makefile" <<'EOF'
.PHONY: lint test

lint:
	@ktn-linter lint ./...

test:
	@go test -v ./...

all: lint test
EOF
        echo -e "${GREEN}✅ Created Makefile${NC}"
    else
        echo -e "${YELLOW}💡 Add these targets to your Makefile:${NC}"
        echo ""
        echo -e "${GREEN}lint:${NC}"
        echo -e "${GREEN}\t@ktn-linter lint ./...${NC}"
        echo ""
    fi
fi

echo ""
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Installation Complete! 🎉           ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}Usage:${NC}"
echo -e "  ${YELLOW}ktn-linter lint ./...${NC}           # Lint your project"
echo -e "  ${YELLOW}ktn-linter lint --help${NC}          # Show help"
echo -e "  ${YELLOW}make lint${NC}                       # If Makefile configured"
echo ""
echo -e "${GREEN}Documentation:${NC}"
echo -e "  ${BLUE}https://github.com/kodflow/ktn-linter${NC}"
echo ""
