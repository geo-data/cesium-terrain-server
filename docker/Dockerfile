##
# Image for the Cesium Terrain Server
#

# Ubuntu 14.04 Trusty Tahr
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
    rsync \
    unzip

# Install Nginx
ADD nginx.list /etc/apt/sources.list.d/nginx.list
RUN apt-key adv --fetch-keys http://nginx.org/keys/nginx_signing.key
RUN apt-get install -y nginx
RUN rm /etc/nginx/sites-enabled/default

# Install Go
ADD install-go.sh /tmp/
RUN /tmp/install-go.sh

# Install the terrain server
ADD local/ /tmp/local/
ADD install-terrain-server.sh /tmp/
RUN /tmp/install-terrain-server.sh

# Install Cesium.js
ADD cesium-version.txt install-cesiumjs.sh /tmp/
RUN /tmp/install-cesiumjs.sh

# Add our filesystem updates
ADD ./root-fs /tmp/root-fs
RUN rsync -a /tmp/root-fs/ /

# Expose the terrain server and apache ports
EXPOSE 8000 80

# Clean up APT when done
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Use baseimage-docker's init system.
CMD ["/sbin/my_init"]
