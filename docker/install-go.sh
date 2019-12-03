#!/usr/bin/env bash
set -e
set -o pipefail

##
# Install the latest stable release of Go
#
# See <https://golang.org/doc/install/source>
#

GO_VERSION=1.13.4
GO_OS=linux
GO_ARCH=amd64
GO_PAYLOAD="go$GO_VERSION.$GO_OS-$GO_ARCH"

cd /usr/local/src
wget -q https://dl.google.com/go/$GO_PAYLOAD.tar.gz
mkdir unpacked
tar -C /usr/local -xzf $GO_PAYLOAD.tar.gz


echo export PATH=$PATH:/usr/local/go/bin >> /etc/bash.bashrc
echo !!!!!!!!!!!!! go installed !!!!!!!!!!!!!
