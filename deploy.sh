#!/bin/sh

export GO111MODULE=on

APP_NAME="cine_stream"
TARGET="release"
ROOT_DIR=`pwd`

rm -rf $TARGET

ENV="prod"

build() {
    mkdir -p "$ROOT_DIR"/$TARGET/bin
    mkdir -p "$ROOT_DIR"/$TARGET/logs
    mkdir -p "$ROOT_DIR"/$TARGET/script
    mkdir -p "$ROOT_DIR"/$TARGET/conf
    chmod 755 "$ROOT_DIR"/$TARGET/bin/

    make linux

    if [ -f "$ROOT_DIR/$APP_NAME"  ]; then
        mv "$ROOT_DIR"/$APP_NAME "$ROOT_DIR"/$TARGET/bin
    fi
    /bin/cp -r "$ROOT_DIR"/conf/"$ENV"/* "$ROOT_DIR"/$TARGET/conf/
    return
}

build

# scp -r release/* xxx@150.158.174.166:/usr/local/app/cine_server