// Implements a server for distributing Cesium terrain tilesets
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/geo-data/cesium-terrain-server/log"
	"github.com/geo-data/cesium-terrain-server/server"
	"github.com/geo-data/cesium-terrain-server/stores"
	"github.com/geo-data/cesium-terrain-server/stores/files"
	"github.com/geo-data/cesium-terrain-server/stores/items/terrain"
	"github.com/geo-data/cesium-terrain-server/stores/tiles"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io"
	l "log"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	/*if len(memcache) > 0 {
		tileStore := tiles.New(NewTileCacheName(), mc.New(memcache))
		stores = append([]*tiles.Store{tileStore}, stores...)
	}*/

	return stores
}

type LogOpt struct {
	Priority log.Priority
}

func NewLogOpt() *LogOpt {
	return &LogOpt{
		Priority: log.LOG_NOTICE,
	}
}

func (this *LogOpt) String() string {
	switch this.Priority {
	case log.LOG_CRIT:
		return "crit"
	case log.LOG_ERR:
		return "err"
	case log.LOG_NOTICE:
		return "notice"
	default:
		return "debug"
	}
}

func (this *LogOpt) Set(level string) error {
	switch level {
	case "crit":
		this.Priority = log.LOG_CRIT
	case "err":
		this.Priority = log.LOG_ERR
	case "notice":
		this.Priority = log.LOG_NOTICE
	case "debug":
		this.Priority = log.LOG_DEBUG
	default:
		return errors.New("choose one of crit, err, notice, debug")
	}
	return nil
}

type multiWriter struct {
	writers []http.ResponseWriter
}

func (t *multiWriter) Header() http.Header {
	for _, w := range t.writers[1:] {
		w.Header()
	}

	return t.writers[0].Header()
}

func (t *multiWriter) WriteHeader(status int) {
	for _, w := range t.writers {
		w.WriteHeader(status)
	}
}

func (t *multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

// MultiWriter is inspired by io.MultiWriter
func MultiWriter(writers ...http.ResponseWriter) http.ResponseWriter {
	w := make([]http.ResponseWriter, len(writers))
	copy(w, writers)
	return &multiWriter{w}
}

type CacheHandler struct {
	mc      *memcache.Client
	handler http.Handler
}

func NewCacheHandler(connstr string, handler http.Handler) http.Handler {
	return &CacheHandler{
		mc:      memcache.New(connstr),
		handler: handler,
	}
}

func (this *CacheHandler) generateKey(r *http.Request) string {
	var u *url.URL
	if referer, ok := r.Header["Referer"]; ok {
		u, _ = url.Parse(referer[0])
	} else {
		// Copy the request URL
		u, _ = url.Parse(r.URL.String())
	}

	return "tiles" + u.RequestURI()
}

func (this *CacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug(fmt.Sprintf("request: %+v", r.Header))
	rec := httptest.NewRecorder()

	// Write to both the recorder and original writer
	tee := MultiWriter(w, rec)
	this.handler.ServeHTTP(tee, r)

	key := this.generateKey(r)
	log.Debug(fmt.Sprintf("setting key: %s", key))
	if err := this.mc.Set(&memcache.Item{Key: key, Value: rec.Body.Bytes()}); err != nil {
		log.Err(err.Error())
	}

	return
}

func main() {
	port := flag.Uint("port", 8000, "the port on which the server listens")
	tilesetRoot := flag.String("dir", ".", "the root directory under which tileset directories reside")
	webRoot := flag.String("web-dir", "", "(optional) the root directory containing static files to be served")
	memcached := flag.String("memcached", "", "(optional) memcached connection string for caching tiles e.g. localhost:11211")
	logging := NewLogOpt()
	flag.Var(logging, "log-level", "level at which logging occurs. One of crit, err, notice, debug")
	flag.Parse()

	// Set the logging
	log.SetLog(l.New(os.Stderr, "", l.LstdFlags), logging.Priority)

	// Generate a list of valid tile stores.
	tileStores := CreateTileStores(*tilesetRoot, *memcached)

	// The tile stores honour the Storer interface, which we also need.
	var stores []stores.Storer
	for _, store := range tileStores {
		stores = append(stores, store)
	}

	r := mux.NewRouter()
	r.HandleFunc("/tilesets/terrain/{tileset}/layer.json", server.LayerHandler(*tilesetRoot, stores))
	r.HandleFunc("/tilesets/terrain/{tileset}/{z:[0-9]+}/{x:[0-9]+}/{y:[0-9]+}.terrain", server.TerrainHandler(tileStores))
	if len(*webRoot) > 0 {
		log.Debug(fmt.Sprintf("serving static resources from %s", *webRoot))
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(*webRoot)))
	}

	handler := server.AddCorsHeader(r)
	if len(*memcached) > 0 {
		log.Debug(fmt.Sprintf("memcached enabled for all resources: %s", *memcached))
		handler = NewCacheHandler(*memcached, handler)
	}

	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, handler))

	log.Notice(fmt.Sprintf("server listening on port %d", *port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Crit(fmt.Sprintf("server failed: %s", err))
		os.Exit(1)
	}
}
