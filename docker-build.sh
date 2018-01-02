#!/bin/sh

docker run --rm -v $PWD:/go/src/imgfit -w /go/src/imgfit golang:1.9 make $1
