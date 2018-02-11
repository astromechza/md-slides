#!/usr/bin/env bash
set -eu

mkdir -p dist/
for OS in "linux windows darwin"; do
  for ARCH in ="amd64 386"; do
    NAME="md-slides.${OS}.${ARCH}"
    if [[ "${OS}" == "windows" ]]; then
      NAME="${NAME}.exe"
    fi
    GOARCH=${ARCH} GOOS=${OS} CGO_ENABLED=0 go build -o "dist/${NAME}" .
  done
done
