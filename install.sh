#!/usr/bin/env sh

# This install script is intended to download and install the latest available
# release of the this project. It will attempt to download the version of the
# binary required by your platform.
#
# Environment variables:
# - INSTALL_DIRECTORY (optional): defaults to current directory
# - OVERRIDE_RELEASE_TAG (optional): defaults to fetching the latest release
# - OVERRIDE_OS (optional): use a specific value for OS (mostly for testing)
# - OVERRIDE_ARCH (optional): use a specific value for ARCH (mostly for testing)
#
# You can install using this script:
# $ curl https://raw.githubusercontent.com/AstromechZA/md-slides/master/install.sh | sh
#
# To adapt this to your own project, replace <project> with your user and project name.

set -e
set -o pipefail

# project urls
PROJECT_URL="https://github.com/AstromechZA/md-slides"
RELEASES_URL="${PROJECT_URL}/releases"
BINARY_NAME="md-slides"

# destination dir can be overriden
INSTALL_DIRECTORY=${INSTALL_DIRECTORY:-.}

downloadJSON() {
    url="$2"

    echo "Fetching $url.."
    if type curl > /dev/null; then
        response=$(curl -s -L -w 'HTTPSTATUS:%{http_code}' -H 'Accept: application/json' "$url")
        body=$(echo "$response" | sed -e 's/HTTPSTATUS\:.*//g')
        code=$(echo "$response" | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    elif type wget > /dev/null; then
        temp=$(mktemp)
        body=$(wget -q --header='Accept: application/json' -O - --server-response --content-on-error "$url" 2> "$temp")
        code=$(awk '/^  HTTP/{print $2}' < "$temp")
    else
        echo "Neither curl nor wget was available to perform http requests."
        exit 1
    fi
    if [ "$code" != 200 ]; then
        echo "Request failed with code $code"
        exit 1
    fi

    eval "$1='$body'"
}

downloadFile() {
    url="$1"
    destination="$2"

    echo "Fetching $url.."
    if type curl > /dev/null; then
        code=$(curl -s -w '%{http_code}' -L "$url" -o "$destination")
    elif type wget > /dev/null; then
        code=$(wget -q -O "$destination" --server-response "$url" 2>&1 | awk '/^  HTTP/{print $2}')
    else
        echo "Neither curl nor wget was available to perform http requests."
        exit 1
    fi

    if [ "$code" != 200 ]; then
        echo "Request failed with code $code"
        exit 1
    fi
}

initArch() {
    ARCH=$(uname -m)
    if [ -n "$OVERRIDE_ARCH" ]; then
        echo "Using OVERRIDE_ARCH"
        ARCH="$OVERRIDE_ARCH"
    fi
    case $ARCH in
        amd64) ARCH="amd64";;
        x86_64) ARCH="amd64";;
        i386) ARCH="386";;
        *) echo "Architecture ${ARCH} is not supported by this installation script"; exit 1;;
    esac
    echo "ARCH = $ARCH"
}

initOS() {
    OS=$(uname | tr '[:upper:]' '[:lower:]')
    if [ -n "$OVERRIDE_OS" ]; then
        echo "Using OVERRIDE_OS"
        OS="$OVERRIDE_OS"
    fi
    case "$OS" in
        darwin) OS='darwin';;
        linux) OS='linux';;
        freebsd) OS='freebsd';;
        mingw*) OS='windows';;
        msys*) OS='windows';;
        windows) OS='windows';;
        *) echo "OS ${OS} is not supported by this installation script"; exit 1;;
    esac
    echo "OS = $OS"
}

# identify platform based on uname output
initArch
initOS

# assemble expected release artifact name
# you will also need to modify this pattern when adapting it to your project
DOWNLOAD_BINARY="${BINARY_NAME}.${OS}.${ARCH}"

# add .exe if on windows
if [ "${OS}" = "windows" ]; then
    DOWNLOAD_BINARY="$DOWNLOAD_BINARY.exe"
fi

# if OVERRIDE_RELEASE_TAG was not provided, assume latest
if [ -z "$OVERRIDE_RELEASE_TAG" ]; then
    downloadJSON LATEST_RELEASE "$RELEASES_URL/latest"
    OVERRIDE_RELEASE_TAG=$(echo "${LATEST_RELEASE}" | tr -s '\n' ' ' | sed 's/.*"tag_name":"//' | sed 's/".*//' )
fi
echo "Release Tag = $OVERRIDE_RELEASE_TAG"

# fetch the real release data to make sure it exists before we attempt a download
downloadJSON RELEASE_DATA "$RELEASES_URL/tag/$OVERRIDE_RELEASE_TAG"

BINARY_URL="$RELEASES_URL/download/$OVERRIDE_RELEASE_TAG/$DOWNLOAD_BINARY"
DOWNLOAD_FILE=$(mktemp)

downloadFile "$BINARY_URL" "$DOWNLOAD_FILE"

echo "Setting executable permissions."
chmod +x "$DOWNLOAD_FILE"


# add .exe if on windows
if [ "${OS}" = "windows" ]; then
    BINARY_NAME="$BINARY_NAME.exe"
fi

echo "Moving executable to $INSTALL_DIRECTORY/$BINARY_NAME"
mv "$DOWNLOAD_FILE" "$INSTALL_DIRECTORY/$BINARY_NAME"
