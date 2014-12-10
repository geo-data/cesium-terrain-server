#!/bin/sh

##
# Install the terrain server
#

export GOPATH=/tmp/go
export PATH=${GOPATH}/bin:$PATH
CTS_DIR=${GOPATH}/src/github.com/geo-data/cesium-terrain-server

# Extract the terrain server code
mkdir -p $CTS_DIR || exit 1
cd $CTS_DIR || exit 1
tar -xzvf /tmp/local/cesium-terrain-server.tar.gz || exit 1

# Build and install the server
make server || exit 1
install --strip ./server /usr/local/bin/terrain-server || exit 1
