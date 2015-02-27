#!/bin/sh

##
# Set up the container based on the environment
#

# If a memcached container has been linked in then use it
if [ -n "${MEMCACHED_PORT_11211_TCP_PORT}" ]
then
    export MEMCACHED="memcached:${MEMCACHED_PORT_11211_TCP_PORT}"
    echo -n $MEMCACHED > /etc/container_environment/MEMCACHED
fi
