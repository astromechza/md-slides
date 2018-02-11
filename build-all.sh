#!/usr/bin/env bash
set -eu

mkdir -p dist/
OSS="linux windows darwin"
ARCHS="amd64 386"
for OS in "${OSS[@]}"; do
  for ARCH in "${ARCHS[@]}"; do
    NAME="md-slides.${OS}.${ARCH}"
    if [[ "${OS}" == "windows" ]]; then
      NAME="${NAME}.exe"
    fi
    GOARCH=${ARCH} GOOS=${OS} CGO_ENABLED=0 go build -o "dist/${NAME}" .
  done
done
