#!/bin/sh

##
# Set up the memcached environment
#
# If the container client has specified a memcached connection we need to make
# Nginx and the terrain server aware of it.  For Nginx this is done by editing
# the config file (as the Nginx config file does not expose environment
# variables).  For the terrain server this is done by setting the $MEMCACHED
# environment variable.
#

# If a memcached container has been linked in then use it
if [ -n "${MEMCACHED_PORT_11211_TCP_PORT}" ]
then
    export MEMCACHED="memcached:${MEMCACHED_PORT_11211_TCP_PORT}"
    echo -n $MEMCACHED > /etc/container_environment/MEMCACHED
fi

# Update the Nginx configuration
if [ -n "$MEMCACHED" ]
then
    perl -p -i -e "s/MEMCACHED/${MEMCACHED}" /etc/nginx/sites-available/cesium-memcached
    ln -s /etc/nginx/sites-available/cesium-memcached /etc/nginx/sites-enabled/
else
    ln -s /etc/nginx/sites-available/cesium-default /etc/nginx/sites-enabled/
fi
