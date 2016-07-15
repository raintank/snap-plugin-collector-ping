#!/bin/bash -e

# Find the directory we exist within
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

# download snap.
go get -d github.com/intelsdi-x/snap
cd $GOPATH/src/github.com/intelsdi-x/snap

# get snap dependencies, will restore sepecific versions defined in Godeps.json
make deps

cd ${DIR}/..
# get all other dependencies
go get ./...