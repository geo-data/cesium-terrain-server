# Cesium Terrain Server

A basic server for serving up filesystem based tilesets representing
[Cesium.js](http://cesiumjs.org/) terrain models.  The resources served up are
intended for use with the
[`CesiumTerrainProvider`](http://cesiumjs.org/Cesium/Build/Documentation/CesiumTerrainProvider.html)
JavaScript class present in the Cesium.js client.

This has specifically been created for easing the development and testing of
terrain tilesets created using the
[Cesium Terrain Builder](https://github.com/geo-data/cesium-terrain-builder)
tools.  On its own it is not recommended for production deployments as no
optimisation (especially caching) has been implemented yet.  Caching, however,
can be added to the stack using third party HTTP caching tools.  The code base
can also provide a useful starting point for developing custom applications.

This project also provides a [Docker](https://www.docker.com/) container to
further simplify deployment of the server and testing of tilesets.  See the
[Docker Registry](https://registry.hub.docker.com/u/geodata/cesium-terrain-server/)
for further details.

## Usage

The terrain server is a self contained binary with the following command line
options:

```sh
$ ./server -help
Usage of server:
  -dir=".": the root directory under which tileset directories reside
  -port=8000: the port on which the server listens
```

Assuming you have a terrain tileset (possibly created with
[`ctb-tile`](https://github.com/geo-data/cesium-terrain-builder#ctb-tile))
located at `/data/tilesets/terrain` which you want served up on port `8080`, you
would run the following command:

```sh
./server -dir /data/tilesets/terrain -port 8080
```

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

## Installation

The server is written in [Go](http://golang.org/) and requires Go to be present
on the system when compiling it from source.  As such, it should run everywhere
that Go does.  Assuming that you have set the
[GOPATH](https://golang.org/cmd/go/#hdr-GOPATH_environment_variable),
installation is a matter of running `go install`:

```sh
go install github.com/geo-data/cesium-terrain-server
```

A program called `cesium-terrain-server` should then be available under your
`GOPATH` (or `GOBIN` location if set).

## Developing

The code has been developed on a Linux platform. After downloading the package
you should be able to run `make server` from the project root to build the
server.

Executing `make docker-local` will create a docker image tagged
`geodata/cesium-terrain-server:local` which when run with a bind mount to the
project source directory is very handy for developing and testing.

## Issues and Contributing

Please report bugs or issues using the
[GitHub issue tracker](https://github.com/geo-data/cesium-terrain-server).

Code and documentation contributions are very welcome, either as GitHub pull
requests or patches.

## License

The [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).

## Contact

Homme Zwaagstra <hrz@geodata.soton.ac.uk>
