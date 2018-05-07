#!/bin/bash

BUILD_DIR=bin
NAME=poorman-sd

mkdir $BUILD_DIR

for GOOS in darwin linux; do
    for GOARCH in amd64; do
        GOOS=$GOOS GOARCH=$GOARCH go build --ldflags="-s -w" -v -o $BUILD_DIR/$NAME-$GOOS-$GOARCH
	# compress binary
	if [ $(which upx) ] ; then
	   upx -5 $BUILD_DIR/$NAME-$GOOS-$GOARCH
	fi
    done
done
