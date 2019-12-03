# Cesium Terrain Server

A basic server for serving up filesystem based tilesets representing
[Cesium.js](http://cesiumjs.org/) terrain models.  The resources served up are
intended for use with the
[`CesiumTerrainProvider`](http://cesiumjs.org/Cesium/Build/Documentation/CesiumTerrainProvider.html)
JavaScript class present in the Cesium.js client.

This has specifically been created for easing the development and testing of
terrain tilesets created using the
[Cesium Terrain Builder](https://github.com/nmccready/cesium-terrain-builder)
tools.

This project also provides a [Docker](https://www.docker.com/) container to
further simplify deployment of the server and testing of tilesets.  See the
[Docker Registry](https://registry.hub.docker.com/u/geodata/cesium-terrain-server/)
for further details.

## Usage

The terrain server is a self contained binary with the following command line
options:

```sh
$ cesium-terrain-server:
  -base-terrain-url="/tilesets": base url prefix under which all tilesets are served
  -cache-limit=1.00MB: the memory size in bytes beyond which resources are not cached. Other memory units can be specified by suffixing the number with kB, MB, GB or TB
  -dir=".": the root directory under which tileset directories reside
  -log-level=notice: level at which logging occurs. One of crit, err, notice, debug
  -memcached="": (optional) memcached connection string for caching tiles e.g. localhost:11211
  -no-request-log=false: do not log client requests for resources
  -port=8000: the port on which the server listens
  -web-dir="": (optional) the root directory containing static files to be served
```

Assume you have the following (small) terrain tileset (possibly created with
[`ctb-tile`](https://github.com/nmccready/cesium-terrain-builder#ctb-tile)):

```
/data/tilesets/terrain/srtm/
├── 0
│   └── 0
│       └── 0.terrain
├── 1
│   └── 1
│       └── 1.terrain
├── 2
│   └── 3
│       └── 3.terrain
└── 3
    └── 7
        └── 6.terrain
```

To serve this tileset on port `8080`, you would run the following command:

```sh
cesium-terrain-server -dir /data/tilesets/terrain -port 8080
```

The tiles would then be available under <http://localhost:8080/tilesets/srtm/>
(e.g. <http://localhost:8080/tilesets/srtm/0/0/0.terrain> for the root tile).
This URL, for instance, is what you would use when configuring
[`CesiumTerrainProvider`](http://cesiumjs.org/Cesium/Build/Documentation/CesiumTerrainProvider.html)
in the Cesium client.

Serving up additional tilesets is simply a matter of adding the tileset as a
subdirectory to `/data/tilesets/terrain/`.  For example, adding a tileset
directory called `lidar` to that location will result in the tileset being
available under <http://localhost:8080/tilesets/lidar/>.

Note that the `-web-dir` option can be used to serve up static assets on the
filesystem in addition to tilesets.  This makes it easy to use the server to
prototype and develop web applications around the terrain data.

### `layer.json`

The `CesiumTerrainProvider` Cesium.js class requires that a `layer.json`
resource is present describing the terrain tileset.  The `ctb-tile` utility does
not create this file.  If a `layer.json` file is present in the root directory
of the tileset then this file will be returned by the server when the client
requests it.  If the file is not found then the server will return a default
resource.

### Root tiles

The Cesium javascript client requires that the two top level tiles representing
zoom level `0` are always present.  These tiles are represented by the
`0/0/0.terrain` and `0/1/0.terrain` resources. When creating tilesets using the
`ctb-tile` utility only one of these tiles will be generated *unless* the source
terrain dataset intersects with the prime meridian.  The terrain server
addresses this issue by serving up a blank terrain tile if a top level tile is
requested which does not also exist on the filesystem.

### Caching tiles with Memcached

The terrain server can use a memcache server to cache tileset data. It is
important to note that the terrain server does not use the cache itself, it only
populates it for each request.  The idea is that a reverse proxy attached to the
memcache (such as Nginx) will first attempt to fulfil a request from the cache
before falling back to the terrain server, which will then update the cache.

Enabling this functionality requires specifying the network address of a
memcached server (including the port) using the `-memcached` option.  E.g. A
memcached server running at `memcache.me.org` on port `11211` can be used as
follows:

```sh
cesium-terrain-server -dir /data/tilesets/terrain -memcached memcache.me.org:11211
```

If present, the terrain server uses the value of the custom `X-Memcache-Key`
header as the memcache key, otherwise it uses the value of the request URI.  A
minimal Nginx configuration setting `X-Memcache-Key` is as follows:

```
server {
    listen 80;

    server_name localhost;

    root /var/www/app;
    index index.html;

    location /tilesets/ {
        set            $memcached_key "tiles$request_uri";
        memcached_pass memcached:11211;
        error_page     404 502 504 = @fallback;
        add_header Access-Control-Allow-Origin "*";

        location ~* \.terrain$ {
            add_header Content-Encoding gzip;
        }
    }

    location @fallback {
        proxy_pass     http://tiles:8000;
        proxy_set_header X-Memcache-Key $memcached_key;
    }
}
```

The `-cache-limit` option can be used in conjunction with the above to change
the memory limit at which resources are considered to large for the cache.

## Installation

The server is written in [Go](http://golang.org/) and requires Go to be present
on the system when compiling it from source.  As such, it should run everywhere
that Go does.  Assuming that you have set the
[GOPATH](https://golang.org/cmd/go/#hdr-GOPATH_environment_variable),
installation is a matter of running `go install`:

```sh
go get github.com/nmccready/cesium-terrain-server/cmd/cesium-terrain-server
```

A program called `cesium-terrain-server` should then be available under your
`GOPATH` (or `GOBIN` location if set).

## Developing

The code has been developed on a Linux platform. After downloading the package
you should be able to run `make` from the project root to build the server,
which will be available as `./bin/cesium-terrain-server`.

Executing `make docker-local` will create a docker image tagged
`geodata/cesium-terrain-server:local` which when run with a bind mount to the
project source directory is very handy for developing and testing.

## Issues and Contributing

Please report bugs or issues using the
[GitHub issue tracker](https://github.com/nmccready/cesium-terrain-server).

Code and documentation contributions are very welcome, either as GitHub pull
requests or patches.

## License

The [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).

## Contact

Homme Zwaagstra <hrz@geodata.soton.ac.uk>
