#!/usr/bin/env sh

set -e

REPO="davenicholson-xyz/wallchemy" 
BINARY_NAME="wallchemy"
VERSION=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize architecture name
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Determine file extension and filename
if [ "$OS" = "windows_nt" ]; then
    EXT="zip"
else
    EXT="tar.gz"
fi

FILE="${BINARY_NAME}-${OS}-${ARCH}-${VERSION}.${EXT}"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILE}"

# Download files
echo "Downloading $FILE..."
curl -LO "$URL"


# Extract and install
if [ "$EXT" = "zip" ]; then
    unzip -o "$FILE"
else
    tar -xzf "$FILE"
fi

chmod +x "$BINARY_NAME"

# Move to /usr/local/bin
INSTALL_PATH="/usr/local/bin/$BINARY_NAME"
echo "Installing to $INSTALL_PATH"
sudo mv "$BINARY_NAME" "$INSTALL_PATH"
rm $FILE

echo "$BINARY_NAME installed successfully!"

