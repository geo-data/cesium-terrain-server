#!/bin/sh

##
# Install the terrain server
#

export GOPATH=/tmp/go
export PATH=${GOPATH}/bin:$PATH

# If a local source code archive does not exist, get it from GitHub.
checkout=`cat /tmp/cts-checkout.txt`
archive="/tmp/local/cesium-terrain-server-${checkout}.tar.gz"
if [ ! -f $archive ]; then
    wget --no-verbose -O $archive "https://github.com/geo-data/cesium-terrain-server/archive/${checkout}.tar.gz" || exit 1
fi

# Set up the source directory
CTS_DIR=${GOPATH}/src/github.com/geo-data/cesium-terrain-server
mkdir -p $CTS_DIR || exit 1
cd $CTS_DIR || exit 1

# Extract the terrain server code
tar --strip-components=1 -xzf $archive || exit 1

# Build and install the server
make server || exit 1
install --strip ./server /usr/local/bin/terrain-server || exit 1
