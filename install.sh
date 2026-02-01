#!/bin/sh
set -e

# Skim installer script
# Usage: curl -sSfL https://raw.githubusercontent.com/Ayushlm10/skim/main/install.sh | sh

REPO="Ayushlm10/skim"
BINARY_NAME="skim"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*)  echo "darwin" ;;
        Linux*)   echo "linux" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *)        echo "unsupported" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *)             echo "unsupported" ;;
    esac
}

# Get latest release version from GitHub
get_latest_version() {
    curl -sSfL "https://api.github.com/repos/${REPO}/releases/latest" | 
        grep '"tag_name":' | 
        sed -E 's/.*"v([^"]+)".*/\1/'
}

main() {
    OS=$(detect_os)
    ARCH=$(detect_arch)

    if [ "$OS" = "unsupported" ]; then
        echo "Error: Unsupported operating system"
        exit 1
    fi

    if [ "$ARCH" = "unsupported" ]; then
        echo "Error: Unsupported architecture"
        exit 1
    fi

    if [ "$OS" = "windows" ] && [ "$ARCH" = "arm64" ]; then
        echo "Error: Windows ARM64 is not supported"
        exit 1
    fi

    echo "Detecting system... ${OS}/${ARCH}"

    VERSION=$(get_latest_version)
    if [ -z "$VERSION" ]; then
        echo "Error: Could not determine latest version"
        exit 1
    fi

    echo "Latest version: v${VERSION}"

    # Construct download URL
    if [ "$OS" = "windows" ]; then
        ARCHIVE_NAME="${BINARY_NAME}_${VERSION}_${OS}_${ARCH}.zip"
    else
        ARCHIVE_NAME="${BINARY_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
    fi

    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${ARCHIVE_NAME}"

    echo "Downloading ${ARCHIVE_NAME}..."

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf ${TMP_DIR}" EXIT

    # Download and extract
    curl -sSfL "${DOWNLOAD_URL}" -o "${TMP_DIR}/${ARCHIVE_NAME}"

    echo "Extracting..."
    if [ "$OS" = "windows" ]; then
        unzip -q "${TMP_DIR}/${ARCHIVE_NAME}" -d "${TMP_DIR}"
    else
        tar -xzf "${TMP_DIR}/${ARCHIVE_NAME}" -C "${TMP_DIR}"
    fi

    # Install binary
    echo "Installing to ${INSTALL_DIR}..."
    if [ -w "$INSTALL_DIR" ]; then
        mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        echo "Need sudo access to install to ${INSTALL_DIR}"
        sudo mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

    echo ""
    echo "Successfully installed ${BINARY_NAME} v${VERSION} to ${INSTALL_DIR}"
    echo ""
    echo "Run 'skim --help' to get started"
}

main
