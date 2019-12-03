#!/usr/bin/env bash
set -e
set -o pipefail

##
# Install the terrain server
#

# If a local source code archive does not exist, get it from GitHub.
archive="/tmp/local/cesium-terrain-server-${FRIENDLY_CHECKOUT}.tar.gz"
if [ ! -f $archive ]; then
	echo !!!! Downloading Archive Local not Found !!!!!
	wget --no-verbose -O $archive "https://github.com/nmccready/cesium-terrain-server/archive/${checkout}.tar.gz"
fi

if [[ -z $GOPATH ]]; then
	GOPATH=$GOROOT
fi
# Set up the source directory
CTS_DIR=$GOPATH/src/github.com/nmccready/cesium-terrain-server
mkdir -p $CTS_DIR
echo made src/github/cst_dir
cd $CTS_DIR
echo "!!! in CTS_DIR !!!"

# Extract the terrain server code
tar --strip-components=1 -xzf $archive
echo "!!! untared archive !!!"

# Build and install the server
make install
# echo "!!! sucessful install !!!"
