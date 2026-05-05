#!/bin/bash

# build-release.sh
# Build and package UAD binaries for release

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get version from argument or git tag
VERSION="${1:-}"
if [ -z "$VERSION" ]; then
    # Try to get version from git tag
    VERSION=$(git describe --tags --exact-match 2>/dev/null || echo "")
    if [ -z "$VERSION" ]; then
        echo -e "${RED}Error: Version not provided and no git tag found${NC}"
        echo "Usage: $0 <version>"
        echo "Example: $0 v0.1.0"
        exit 1
    fi
fi

# Remove 'v' prefix if present
VERSION_NUMBER="${VERSION#v}"

echo -e "${GREEN}Building UAD Language Release: ${VERSION}${NC}"
echo ""

# Create release directory
RELEASE_DIR="release"
mkdir -p "$RELEASE_DIR"

# Get current directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"
cd "$PROJECT_ROOT"

# Build function
build_for_platform() {
    local GOOS=$1
    local GOARCH=$2
    local PLATFORM_NAME=$3
    
    echo -e "${YELLOW}Building for ${PLATFORM_NAME}...${NC}"
    
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    
    # Clean previous builds
    rm -rf bin/*
    
    # Build all binaries
    echo "  Building uadc..."
    go build -ldflags="-s -w -X main.Version=${VERSION}" -o bin/uadc ./cmd/uadc
    
    echo "  Building uadi..."
    go build -ldflags="-s -w -X main.Version=${VERSION}" -o bin/uadi ./cmd/uadi
    
    echo "  Building uadvm..."
    go build -ldflags="-s -w -X main.Version=${VERSION}" -o bin/uadvm ./cmd/uadvm
    
    echo "  Building uadrepl..."
    go build -ldflags="-s -w -X main.Version=${VERSION}" -o bin/uadrepl ./cmd/uadrepl
    
    echo "  Building uad-runner..."
    go build -ldflags="-s -w -X main.Version=${VERSION}" -o bin/uad-runner ./cmd/uad-runner
    
    # Create platform directory
    PLATFORM_DIR="${RELEASE_DIR}/uad-${VERSION_NUMBER}-${PLATFORM_NAME}"
    mkdir -p "$PLATFORM_DIR"
    
    # Copy binaries
    if [ "$GOOS" = "windows" ]; then
        cp bin/uadc.exe "$PLATFORM_DIR/" 2>/dev/null || true
        cp bin/uadi.exe "$PLATFORM_DIR/" 2>/dev/null || true
        cp bin/uadvm.exe "$PLATFORM_DIR/" 2>/dev/null || true
        cp bin/uadrepl.exe "$PLATFORM_DIR/" 2>/dev/null || true
        cp bin/uad-runner.exe "$PLATFORM_DIR/" 2>/dev/null || true
    else
        cp bin/uadc "$PLATFORM_DIR/" 2>/dev/null || true
        cp bin/uadi "$PLATFORM_DIR/" 2>/dev/null || true
        cp bin/uadvm "$PLATFORM_DIR/" 2>/dev/null || true
        cp bin/uadrepl "$PLATFORM_DIR/" 2>/dev/null || true
        cp bin/uad-runner "$PLATFORM_DIR/" 2>/dev/null || true
    fi
    
    # Copy documentation
    cp README.md "$PLATFORM_DIR/" 2>/dev/null || true
    [ -f LICENSE ] && cp LICENSE "$PLATFORM_DIR/" 2>/dev/null || true
    
    # Create archive
    cd "$RELEASE_DIR"
    if [ "$GOOS" = "windows" ]; then
        ARCHIVE_NAME="uad-${VERSION_NUMBER}-${PLATFORM_NAME}.zip"
        if command -v zip >/dev/null 2>&1; then
            zip -r "$ARCHIVE_NAME" "uad-${VERSION_NUMBER}-${PLATFORM_NAME}" > /dev/null
        else
            echo -e "${RED}Warning: zip command not found, skipping archive creation${NC}"
        fi
    else
        ARCHIVE_NAME="uad-${VERSION_NUMBER}-${PLATFORM_NAME}.tar.gz"
        tar -czf "$ARCHIVE_NAME" "uad-${VERSION_NUMBER}-${PLATFORM_NAME}"
    fi
    
    # Generate checksum
    if [ -f "$ARCHIVE_NAME" ]; then
        if command -v shasum >/dev/null 2>&1; then
            shasum -a 256 "$ARCHIVE_NAME" > "${ARCHIVE_NAME}.sha256"
        elif command -v sha256sum >/dev/null 2>&1; then
            sha256sum "$ARCHIVE_NAME" > "${ARCHIVE_NAME}.sha256"
        fi
        echo -e "${GREEN}  Created ${ARCHIVE_NAME}${NC}"
        echo -e "${GREEN}  Created ${ARCHIVE_NAME}.sha256${NC}"
    fi
    
    cd "$PROJECT_ROOT"
    echo ""
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

echo "Go version: $(go version)"
echo ""

# Detect current platform and build for it
CURRENT_OS=$(go env GOOS)
CURRENT_ARCH=$(go env GOARCH)

echo -e "${GREEN}Building for current platform: ${CURRENT_OS}-${CURRENT_ARCH}${NC}"
echo ""

# Build for current platform
if [ "$CURRENT_OS" = "windows" ]; then
    PLATFORM_NAME="windows-${CURRENT_ARCH}"
else
    PLATFORM_NAME="${CURRENT_OS}-${CURRENT_ARCH}"
fi

build_for_platform "$CURRENT_OS" "$CURRENT_ARCH" "$PLATFORM_NAME"

echo -e "${GREEN}Release build complete!${NC}"
echo ""
echo "Release files are in: ${RELEASE_DIR}/"
echo ""
echo "To build for other platforms, set GOOS and GOARCH environment variables:"
echo "  GOOS=linux GOARCH=amd64 $0 ${VERSION}"
echo "  GOOS=darwin GOARCH=arm64 $0 ${VERSION}"
echo "  GOOS=windows GOARCH=amd64 $0 ${VERSION}"


