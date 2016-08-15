#!/bin/bash

die() {
	echo $1
	exit 1
}

file ../superside | grep "ELF.*LSB" || die "../superside is missing or not a Linux binary"
test -f ../superside.toml || die "Missing superside.toml file!"

cd ../public && npm install
cd .. && docker build -f docker/Dockerfile -t superside . || die "Failed to build"
