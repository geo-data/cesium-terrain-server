// Implements a server for distributing Cesium terrain tilesets
package main

import (
	"flag"
	"fmt"
	"github.com/geo-data/cesium-terrain-server/server"
	"github.com/geo-data/cesium-terrain-server/stores"
	"github.com/geo-data/cesium-terrain-server/stores/files"
	"github.com/geo-data/cesium-terrain-server/stores/items/terrain"
	mc "github.com/geo-data/cesium-terrain-server/stores/memcache"
	"github.com/geo-data/cesium-terrain-server/stores/tiles"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type TileFileName struct {
}

func NewTileFileName() tiles.Namer {
	return &TileFileName{}
}

func (this *TileFileName) TileName(tileset string, tile *terrain.Terrain) string {
	return filepath.Join(
		tileset,
		strconv.FormatUint(tile.Z, 10),
		strconv.FormatUint(tile.X, 10),
		strconv.FormatUint(tile.Y, 10)+".terrain")
}

type TileCacheName struct {
}

func NewTileCacheName() tiles.Namer {
	return &TileCacheName{}
}

func (this *TileCacheName) TileName(tileset string, tile *terrain.Terrain) string {
	return fmt.Sprintf("%s-%d-%d-%d", tileset, tile.Z, tile.X, tile.Y)
}

func CreateTileStores(tilesetRoot, memcache string) []*tiles.Store {
	// There will always be a base file system store
	stores := []*tiles.Store{
		tiles.New(NewTileFileName(), files.New(tilesetRoot)),
	}

	// If a memcache server has been specified, prepend it to the list of stores.
	if len(memcache) > 0 {
		tileStore := tiles.New(NewTileCacheName(), mc.New(memcache))
		stores = append([]*tiles.Store{tileStore}, stores...)
	}

	return stores
}

func main() {
	port := flag.Uint("port", 8000, "the port on which the server listens")
	tilesetRoot := flag.String("dir", ".", "the root directory under which tileset directories reside")
	memcache := flag.String("memcache", "", "memcache connection string for caching tiles e.g. localhost:11211")
	flag.Parse()

	// Generate a list of valid tile stores.
	tileStores := CreateTileStores(*tilesetRoot, *memcache)

	// The tile stores honour the Storer interface, which we also need.
	var stores []stores.Storer
	for _, store := range tileStores {
		stores = append(stores, store)
	}

	r := mux.NewRouter()
	r.HandleFunc("/tilesets/{tileset}/layer.json", server.LayerHandler(*tilesetRoot, stores))
	r.HandleFunc("/tilesets/{tileset}/{z:[0-9]+}/{x:[0-9]+}/{y:[0-9]+}.terrain", server.TerrainHandler(tileStores))

	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, server.AddCorsHeader(r)))

	log.Println("Terrain server listening on port", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
