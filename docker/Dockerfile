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

# Install Apache
RUN apt-get install -y \
    apache2 \
    apache2-utils \
    unzip
RUN mkdir /var/lock/apache2 /var/run/apache2
RUN a2enmod proxy_http
RUN a2dissite 000-default

# Install Go
ADD install-go.sh /tmp/
RUN /tmp/install-go.sh

# Install the terrain server
ADD local/ /tmp/local/
ADD install-terrain-server.sh /tmp/
RUN /tmp/install-terrain-server.sh

# Install Cesium.js
ADD install-cesiumjs.sh /tmp/
RUN /tmp/install-cesiumjs.sh

# Add our filesystem updates
ADD ./root-fs /tmp/root-fs
RUN rsync -a /tmp/root-fs/ /

# Enable the application
RUN a2enconf fqdn
RUN a2ensite terrain-server

# Expose the terrain server and apache ports
EXPOSE 8000 80

# Clean up APT when done
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Use baseimage-docker's init system.
CMD ["/sbin/my_init"]
