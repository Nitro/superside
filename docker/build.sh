#!/bin/bash

die() {
	echo $1
	exit 1
}

file ../superside | grep "ELF.*LSB" || die "../superside is missing or not a Linux binary"
test -f ../superside.toml || die "Missing superside.toml file!"

cd ../public && npm install
cd .. && docker build -f docker/Dockerfile -t superside . || die "Failed to build"
docker tag superside gonitro/superside:latest
docker push gonitro/superside:latest

TAG=`git rev-parse --short HEAD`
docker tag superside gonitro/superside:${TAG}
docker push gonitro/superside:${TAG}
