#!/bin/bash

VERSION="0.0.1"
NAME=elm-recompile

echo "Building $NAME version $VERSION"

mkdir -p pkg

build() {
  echo -n "=> $1-$2: "
  GOOS=$1 GOARCH=$2 go build -o pkg/$NAME-$1-$2 -ldflags "-X main.version=$VERSION -X main.gitHash=`git rev-parse HEAD`" ./main.go
  du -h pkg/$NAME-$1-$2
}

build "darwin" "amd64"
