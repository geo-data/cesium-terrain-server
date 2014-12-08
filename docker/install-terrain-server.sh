#!/bin/sh

##
# Install the terrain server
#

mkdir /tmp/local/cesium-terrain-server || exit 1
cd /tmp/local/cesium-terrain-server || exit 1
tar -xzvf ../cesium-terrain-server.tar.gz || exit 1
make server || exit 1
install --strip ./server /usr/local/bin/terrain-server || exit 1
