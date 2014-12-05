##
# Image for the Cesium Terrain Server
#

# Ubuntu 14.04
FROM phusion/baseimage:0.9.15

# Set correct environment variables.
ENV HOME /root

# Regenerate SSH host keys. baseimage-docker does not contain any, so you
# have to do that yourself. You may also comment out this instruction; the
# init system will auto-generate one during boot.
RUN /etc/my_init.d/00_regen_ssh_host_keys.sh

# Update the location of the apt sources
RUN apt-get update -y

# Install dependencies
RUN apt-get install -y \
    wget \
    build-essential \
    git \
    mercurial \
    rsync

# Install Go
ADD install-go.sh /tmp/
RUN /tmp/install-go.sh

# Install the terrain server
ADD cesium-terrain-server.tar.gz /tmp/cesium-terrain-server/
RUN cd /tmp/cesium-terrain-server && \
    make server && \
    install --strip ./server /usr/local/bin/terrain-server

# Add our filesystem updates
ADD ./root-fs /tmp/root-fs
RUN rsync -a /tmp/root-fs/ /

# Expose the terrain server port
EXPOSE 8000

# Clean up APT when done
#RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Use baseimage-docker's init system.
CMD ["/sbin/my_init"]
