#!/bin/bash

set -eu

NAME=$1
VERSION=$2
LDFLAGS=$3
MAIN_FILES=${@:4:($#-2)}

for os in darwin linux windows
do
	for arch in amd64 386
	do
		binname=$NAME
		if [ "$os" = "windows" ]; then
			binname=$NAME.exe
		fi
		binpath=dist/$VERSION/${NAME}_${os}_$arch/$binname
		GOOS=$os GOARCH=$arch CGO_ENABLED=0 \
			go build -a -tags netgo -installsuffix netgo "$LDFLAGS" \
			-o $binpath $MAIN_FILES
	done
done
