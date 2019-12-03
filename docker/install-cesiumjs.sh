#!/usr/bin/env bash
set -e
set -o pipefail

##
# Install the latest stable release of CesiumJS
#

mkdir -p /tmp/cesium /var/www/cesium || exit 1
cd /tmp/cesium || exit 1

# Get Cesium if we need to
if [ ! -f /tmp/local/Cesium-${CESIUM_VERSION}.zip ]; then
    wget --no-verbose --directory-prefix=/tmp/local https://cesiumjs.org/releases/Cesium-${CESIUM_VERSION}.zip || exit 1
fi

unzip -q /tmp/local/Cesium-${CESIUM_VERSION}.zip || exit 1
mv Apps ThirdParty Build /var/www/cesium/ || exit 1
