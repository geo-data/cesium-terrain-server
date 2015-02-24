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

if [ ! -f /tmp/local/cesium-terrain-server.tar.gz ]; then
    git clone https://github.com/geo-data/cesium-terrain-server.git . || exit 1
else
    tar -xzvf /tmp/local/cesium-terrain-server.tar.gz || exit 1
fi

# Build and install the server
make bin/cesium-terrain-server || exit 1
install --strip bin/cesium-terrain-server /usr/local/bin/terrain-server || exit 1
