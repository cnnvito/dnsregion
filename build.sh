#!/bin/bash

cd $(dirname $(readlink -f $0))

export CGO_ENABLED=0
export GO111MODULE=on

[ ! -d bin/ ] && mkdir bin/

go build -trimpath -ldflags "-s -w" -o bin/dnsregion cmd/*
