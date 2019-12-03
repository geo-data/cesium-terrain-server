// Implements a server for distributing Cesium terrain tilesets
package main

import (
	"flag"
	"fmt"
	myhandlers "github.com/nmccready/cesium-terrain-server/handlers"
	"github.com/nmccready/cesium-terrain-server/log"
	"github.com/nmccready/cesium-terrain-server/stores/fs"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	l "log"
	"net/http"
	"os"
)

func main() {
	port := flag.Uint("port", 8000, "the port on which the server listens")
	tilesetRoot := flag.String("dir", ".", "the root directory under which tileset directories reside")
	webRoot := flag.String("web-dir", "", "(optional) the root directory containing static files to be served")
	memcached := flag.String("memcached", "", "(optional) memcached connection string for caching tiles e.g. localhost:11211")
	baseTerrainUrl := flag.String("base-terrain-url", "/tilesets", "base url prefix under which all tilesets are served")
	noRequestLog := flag.Bool("no-request-log", false, "do not log client requests for resources")
	logging := NewLogOpt()
	flag.Var(logging, "log-level", "level at which logging occurs. One of crit, err, notice, debug")
	limit := NewLimitOpt()
	limit.Set("1MB")
	flag.Var(limit, "cache-limit", `the memory size in bytes beyond which resources are not cached. Other memory units can be specified by suffixing the number with kB, MB, GB or TB`)
	flag.Parse()

	// Set the logging
	log.SetLog(l.New(os.Stderr, "", l.LstdFlags), logging.Priority)

	// Get the tileset store
	store := fs.New(*tilesetRoot)

	r := mux.NewRouter()
	r.HandleFunc(*baseTerrainUrl+"/{tileset}/layer.json", myhandlers.LayerHandler(store))
	r.HandleFunc(*baseTerrainUrl+"/{tileset}/{z:[0-9]+}/{x:[0-9]+}/{y:[0-9]+}.terrain", myhandlers.TerrainHandler(store))
	if len(*webRoot) > 0 {
		log.Debug(fmt.Sprintf("serving static resources from %s", *webRoot))
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(*webRoot)))
	}

	handler := myhandlers.AddCorsHeader(r)
	if len(*memcached) > 0 {
		log.Debug(fmt.Sprintf("memcached enabled for all resources: %s", *memcached))
		handler = myhandlers.NewCache(*memcached, handler, limit.Value, myhandlers.NewLimit)
	}

	if *noRequestLog == false {
		handler = handlers.CombinedLoggingHandler(os.Stdout, handler)
	}

	http.Handle("/", handler)

	log.Notice(fmt.Sprintf("server listening on port %d", *port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Crit(fmt.Sprintf("server failed: %s", err))
		os.Exit(1)
	}
}
