#!/bin/bash

#apt-get update
#apt-get -y install upx

go install github.com/gobuffalo/packr/packr

pushd ../combat
go build -ldflags="-s -w"
upx --brute combat
popd

pushd ../combat-client
go build -ldflags="-s -w"
upx --brute combat-client
popd

pushd ../combat-worker
go build -ldflags="-s -w"
upx --brute combat-worker
popd

packr build -ldflags="-s -w"
upx --brute combat-server
