#!/bin/bash

# Installation script for ktn-linter from GitHub releases
# Downloads the latest release binary and installs it in builds/

set -e

REPO="kodflow/ktn-linter"
INSTALL_DIR="builds"
BINARY_NAME="ktn-linter"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🔧 Installing ktn-linter from GitHub releases...${NC}"

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

echo -e "${YELLOW}📦 Platform detected: ${OS}-${ARCH}${NC}"

# Get latest release URL
echo -e "${YELLOW}🔍 Fetching latest release...${NC}"
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest")
DOWNLOAD_URL=$(echo "$LATEST_RELEASE" | grep "browser_download_url.*${OS}-${ARCH}" | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo -e "${RED}❌ No release found for ${OS}-${ARCH}${NC}"
    echo -e "${YELLOW}💡 Building from source instead...${NC}"
    make build
    exit 0
fi

VERSION=$(echo "$LATEST_RELEASE" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
echo -e "${GREEN}✅ Found version: ${VERSION}${NC}"

# Create install directory
mkdir -p "$INSTALL_DIR"

# Download binary
TEMP_FILE=$(mktemp)
echo -e "${YELLOW}⬇️  Downloading from: ${DOWNLOAD_URL}${NC}"
curl -L -o "$TEMP_FILE" "$DOWNLOAD_URL"

# Install binary
chmod +x "$TEMP_FILE"
mv "$TEMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"

echo -e "${GREEN}✅ ktn-linter installed successfully in ${INSTALL_DIR}/${BINARY_NAME}${NC}"

# Verify installation
if [ -x "${INSTALL_DIR}/${BINARY_NAME}" ]; then
    VERSION_OUTPUT=$("${INSTALL_DIR}/${BINARY_NAME}" --version 2>&1 || echo "version check failed")
    echo -e "${GREEN}✅ Installation verified: ${VERSION_OUTPUT}${NC}"
else
    echo -e "${RED}❌ Installation verification failed${NC}"
    exit 1
fi

echo -e "${GREEN}🎉 Installation complete!${NC}"
