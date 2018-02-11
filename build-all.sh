#!/usr/bin/env bash
set -eu

VERSION=$(git describe --tags --dirty)
COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null)
DATE=$(date "+%Y-%m-%d")

mkdir -p dist/
OSS="linux windows darwin"
ARCHS="amd64 386"
for OS in ${OSS[@]}; do
  for ARCH in ${ARCHS[@]}; do
    NAME="md-slides.${OS}.${ARCH}"
    if [[ "${OS}" == "windows" ]]; then
      NAME="${NAME}.exe"
    fi
    echo "Building ${OS} ${ARCH}"
    GOARCH=${ARCH} GOOS=${OS} CGO_ENABLED=0 go build \
        -o "dist/${NAME}" \
        -ldflags "-X main.commitHash=${COMMIT_HASH} -X main.buildDate=${DATE} -X main.gitVersion=${VERSION}" \
        .
  done
done
