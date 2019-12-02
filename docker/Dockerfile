##
# Image for the Cesium Terrain Server
#

FROM phusion/baseimage:0.9.15
# FROM debian:jessie-slim

ARG FRIENDLY_CHECKOUT
ARG CESIUM_VERSION

# Set correct environment variables.
ENV HOME /root
ENV FRIENDLY_CHECKOUT $FRIENDLY_CHECKOUT
ENV CESIUM_VERSION $CESIUM_VERSION

# Regenerate SSH host keys. baseimage-docker does not contain any, so you
# have to do that yourself. You may also comment out this instruction; the
# init system will auto-generate one during boot.
RUN /etc/my_init.d/00_regen_ssh_host_keys.sh

# Update the location of the apt sources
RUN apt-get update -y

# Install dependencies
RUN apt-get install -y \
		vim \
		curl \
    sed \
    wget \
    bash \
    build-essential \
    git \
    rsync \
    unzip

# RUN apt-get upgrade -y

# default all users to bash
RUN sed -i 's/\/bin\/sh/\/bin\/bash/g' /etc/passwd

# Install Go
ADD install-go.sh /tmp/
RUN /tmp/install-go.sh

# Set the Go workspace
RUN echo 'export PATH="/usr/local/go/bin:$PATH"' >> /etc/bash.bashrc
RUN echo 'export GO111MODULE=on' >> /etc/bash.bashrc

ENV GOROOT=/usr/local/go
ENV GOBIN=/usr/local/go/bin

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

# Expose the terrain server
EXPOSE 8000

# Clean up APT when done
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Use baseimage-docker's init system.
CMD ["/sbin/my_init"]
