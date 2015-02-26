// Implements a server for distributing Cesium terrain tilesets
package main

import (
	"errors"
	"flag"
	"fmt"
	myhandlers "github.com/geo-data/cesium-terrain-server/handlers"
	"github.com/geo-data/cesium-terrain-server/log"
	"github.com/geo-data/cesium-terrain-server/stores/fs"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	l "log"
	"net/http"
	"os"
)

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

func main() {
	port := flag.Uint("port", 8000, "the port on which the server listens")
	tilesetRoot := flag.String("dir", ".", "the root directory under which tileset directories reside")
	webRoot := flag.String("web-dir", "", "(optional) the root directory containing static files to be served")
	memcached := flag.String("memcached", "", "(optional) memcached connection string for caching tiles e.g. localhost:11211")
	baseTerrainUrl := flag.String("base-terrain-url", "/tilesets", "base url prefix under which all tilesets are served")
	noRequestLog := flag.Bool("no-request-log", false, "do not log client requests for resources")
	logging := NewLogOpt()
	flag.Var(logging, "log-level", "level at which logging occurs. One of crit, err, notice, debug")
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
		handler = myhandlers.NewCache(*memcached, handler)
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
