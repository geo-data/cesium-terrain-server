#!/bin/sh

##
# Install the latest stable release of Go
#
# See <https://golang.org/doc/install/source>
#

cd /usr/local/src || exit 1
hg clone -u release https://code.google.com/p/go || exit 1
cd go/src || exit 1
./all.bash || exit 1
ln -s /usr/local/src/go/bin/go* /usr/local/bin/ || exit 1
