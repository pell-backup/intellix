#!/bin/bash

set -e

VERSION=${VERSION:-$(grep 'github.com/IntelliXLabs/iwasm' go.mod | grep -o 'v[0-9.]\+')}
PLATFORM=${PLATFORM:-$(uname -s | tr '[:upper:]' '[:lower:]')}
ARCH=${ARCH:-$(uname -m)}
PROJECT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)

# Map architecture names
case ${ARCH} in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
esac

## Download the lib
DOWNLOAD_URL="https://github.com/IntelliXLabs/iwasm/releases/download/${VERSION}/libruntime-${PLATFORM}-${ARCH}.tar.gz"
INSTALL_DIR="${PROJECT_DIR}/lib"
mkdir -p "${INSTALL_DIR}"

echo "Downloading libruntime from ${DOWNLOAD_URL}"
HTTP_STATUS=$(curl -L \
    -w "%{http_code}" \
    -o "${INSTALL_DIR}/libruntime.tar.gz" \
    "${DOWNLOAD_URL}")

if [ "$HTTP_STATUS" -ne 200 ]; then
    echo "Error: Failed to download library. HTTP status: ${HTTP_STATUS}"
    exit 1
fi

echo "Extracting library..."
tar -xzf "${INSTALL_DIR}/libruntime.tar.gz" -C "${INSTALL_DIR}"

echo "Installation completed successfully"
