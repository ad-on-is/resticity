#!/bin/sh
VERSION=$2

if [ -z "$VERSION" ]; then
  VERSION=$(git describe --tags)
fi

echo $VERSION
exit
BUILD=$(date +%FT%T%z)
LD_FLAGS="-X main.Version=$VERSION -X main.Build=$BUILD"
wails=$(which wails)

case $1 in
  "desktop")
    CGO_ENABLED=0 $wails build -ldflags="$LD_FLAGS"
    ;;
  "server")
    CGO_ENABLED=0 go build -ldflags="$LD_FLAGS" ./cmd/server
    ;;
  "frontend")
    npm i -g pnpm
    cd frontend && pnpm install && pnpm build
    ;;
  "dev")
    $wails dev -ldflags="$LD_FLAGS" -loglevel "Error"
    ;;
  *)
    ;;
esac
