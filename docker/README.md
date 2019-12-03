# Cesium Terrain Server

This provides a Docker container for running
[Cesium Terrain Server](https://github.com/nmccready/cesium-terrain-server). It
is capable of serving custom terrain tilesets with the option of also serving
[Cesium.js](http://cesiumjs.org/) client side assets: it comes bundled with
Cesium.js. This makes it a suitable platform for using in the development of
Cesium.js web applications.

## Usage

```sh
docker run -p 8080:8000 -v /data/docker/tilesets/terrain:/data/tilesets/terrain \
    nmccready/cesium-terrain-server
```

Note that if the `/data/tilesets/terrain` directory is not present on the
container the terrain server will fail to run, as this is where it expects to
find subdirectories representing terrain tilesets.

Running the previous command will serve up the Cesium resources on
<http://localhost:8080/>. For instance you will be able to access the Cesium
Hello World application at <http://localhost:8080/Apps/HelloWorld.html>.

The only change in the container to the stock Cesium.js download is the addition
of a top level `index.html` (this will be the resource returned by
<http://localhost:8080/>). `index.html` is a very minimally modified version of
`Apps/HelloWorld.html` provided by Cesium which additionally loads terrain
pointed to by the the url `/tilesets/test`. In the above example this
requires the host directory `/data/docker/tilesets/terrain/test` to contain a
terrain tileset, which will in turn expose the tileset to
`/data/tilesets/terrain/test` in the container.

### Environment options

A number of environment variables can be set on the container using `docker run --env` which affect its behaviour:

- Setting `SERVE_STATIC=0` will prevent the terrain server for serving up any of
  the static Cesium.js assets. Useful if you only want to serve terrain
  tilesets.

- Setting `WEB_DIR` to a directory will specify the directory from which to
  serve up the static assets. By default this is `/var/www/cesium`. Note also
  setting `SERVE_STATIC=0` will mean that no assets are served, no matter where
  `WEB_DIR` points to.

- Setting `BASE_TERRAIN_PATH` affects the base path under which terrain tiles
  are served. By default this is `/tilesets` but can be set to
  e.g. `BASE_TERRAIN_PATH=/tiles/terrain`.

- Setting `LOG_LEVEL` to one of `crit`, `err`, `notice`, `debug` will determine
  what is logged. `LOG_LEVEL=debug` is useful when developing, especially
  checking what keys are being set in Memcached.

- Setting `LOG_REQUESTS=0` will prevent the logging of HTTP requests.

- Setting `MEMCACHED` enables caching tiles with Memcached. See the relevant
  section below for details.

## Creating and serving tilesets

This container has been designed to be work with the
[Cesium Terrain Builder](https://registry.hub.docker.com/u/homme/cesium-terrain-builder/)
docker image to simplify creating and viewing terrain tilesets. Assume you have
the following folder structure available on your host system:

```
/data/docker/
├── rasters
│   └── DEM.tif
└── tilesets
    └── terrain
```

where `/data/docker/rasters/DEM.tif` contains the height data in GeoTiff format
that you would like to view in Cesium.

First you would use the `ctb-tile` command to create the terrain tileset:

```sh
docker run -v /data/docker:/data homme/cesium-terrain-builder \
    ctb-tile --output-dir /data/tilesets/terrain/test /data/rasters/DEM.tif
```

Then you would serve up the terrain data using Cesium as per the general usage
instructions above:

```sh
docker run -p 8080:80 -v /data/docker/tilesets/terrain:/data/tilesets/terrain \
    nmccready/cesium-terrain-server
```

The terrain data should now be visible at <http://localhost:8080/>.

## Caching tiles with Memcached

The terrain server running within the container can be configured to use a
memcache server to cache tileset data and increase performance. Note that the
terrain server does not use the cache itself, it only populates for each
request. The idea is that a reverse proxy attached to the memcache (such as
Nginx) will first attempt to fulfil a request from the cache before falling back
to the terrain server, which will then update the cache.

Caching is enabled by either specifying a memcached container to link to or
setting the `MEMCACHED` environment variable.

### Linking

Any container running memcached on port 11211 and linked with the alias
`memcached` will be used. E.g. assume the following memcached container is
running:

```sh
docker run --name memcache -d memcached
```

This can then be used by a terrain server image:

```sh
docker run --name terrain -d --link memcache:memcached nmccready/cesium-terrain-server
```

### `MEMCACHED` Environment variable

A memcached server that is not linked can be still used by setting the container
`MEMCACHED` environment variable to point to the memcached network address e.g.

```sh
docker run --name terrain -d --env MEMCACHED=memcache.me.org:11211 nmccready/cesium-terrain-server
```

Linking takes precedence over setting `MEMCACHED`.

## Developing Cesium applications

You can use the container as a base for developing bespoke Cesium applications
with custom terrain data. General workflow would be to create tilesets as
described in the previous section. You would then need to edit
`/var/www/cesium/index.html` in the container to suit your taste. You may also
want to customise the terrain server itself. See the
[Cesium Terrain Server](https://github.com/nmccready/cesium-terrain-server)
project repository for further details on this.

Note that the terrain server will serve up any terrain tilesets present as
subdirectories of `/data/tilesets/terrain` in the container. You are therefore
not limited to the default `test` directory, and you can serve multiple terrain
tilesets at once. To make tilesets other than the default `test` available to
the Cesium client, however, you will need to edit or replace the `index.html`
file to appropriately reference these alternative tileset resources.

You can edit `/var/www/cesium/index.html` directly in the container, but the
recommended approach would be to use the container as a base for your own
application. To do this:

- Clone or download the
  [Cesium Terrain Server](https://github.com/nmccready/cesium-terrain-server)
  repository.
- Edit `docker/root-fs/var/www/cesium/index.html` to suit.
- Build the container defined by the context in `docker/`.

The recommended way to build the container on GNU/Linux distributions is to take
advantage of the `Makefile` in the project root: running `make docker-local`
will create a docker image tagged `nmccready/cesium-terrain-server:local`. This
image, when run with a bind mount to the project root directory, is very handy
for developing and testing.

## Logging

All output from the terrain server is logged under `/var/log/terrain-server`.
This log is managed by [svlogd](http://smarden.org/runit/svlogd.8.html).

## Issues and Contributing

Please report bugs or issues using the
[GitHub issue tracker](https://github.com/nmccready/cesium-terrain-server).

Code and documentation contributions are very welcome, either as GitHub pull
requests or patches.

## License

The [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).

## Contact

Homme Zwaagstra <hrz@geodata.soton.ac.uk>
