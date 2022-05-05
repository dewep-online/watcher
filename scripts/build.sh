#!/bin/bash

cd $PWD
rm -rf $PWD/build/bin/watcher_$1
go generate ./...

# sudo apt install gcc-aarch64-linux-gnu gcc-multilib
case $1 in
arm64)
    GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc \
        go build -a -o $PWD/build/bin/watcher_$1 $PWD/cmd/watcher
  ;;
*)
  GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=$1 \
        go build -o $PWD/build/bin/watcher_$1 $PWD/cmd/watcher
  ;;
esac