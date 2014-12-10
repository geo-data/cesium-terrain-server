// Implements a server for distributing Cesium terrain tilesets
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/geo-data/cesium-terrain-server/assets"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Representation of a terrain tile. This includes the x, y, z coordinate and
// the byte sequence of the tile itself. Note that terrain tiles are gzipped.
type Terrain struct {
	x, y, z uint64
	body    []byte
}

// Load a terrain tile on disk into the Terrain structure
func (self *Terrain) loadFromFs(tilesetRoot string) error {
	filename := filepath.Join(
		tilesetRoot,
		strconv.FormatUint(self.z, 10),
		strconv.FormatUint(self.x, 10),
		strconv.FormatUint(self.y, 10)+".terrain")
	body, err := ioutil.ReadFile(filename)
	if err == nil {
		self.body = body
	}
	return err
}

// Parse x, y, z string coordinates and assign them to the Terrain
func (self *Terrain) parseCoord(x, y, z string) error {
	xi, err := strconv.ParseUint(x, 10, 64)
	if err != nil {
		return err
	}

	yi, err := strconv.ParseUint(y, 10, 64)
	if err != nil {
		return err
	}

	zi, err := strconv.ParseUint(z, 10, 64)
	if err != nil {
		return err
	}

	self.x = xi
	self.y = yi
	self.z = zi

	return nil
}

// An HTTP handler which returns a terrain tile resource
func terrainHandler(tilesetRoot string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var t Terrain

		// get the tile coordinate from the URL
		vars := mux.Vars(r)
		err := t.parseCoord(vars["x"], vars["y"], vars["z"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// try and read the file from disk
		rootDir := filepath.Join(tilesetRoot, vars["tileset"])
		err = t.loadFromFs(rootDir)
		if err != nil {
			if os.IsNotExist(err) {
				if vars["z"] == "0" && vars["y"] == "0" && (vars["x"] == "0" || vars["x"] == "1") {
					// serve up a blank tile as it is a missing root tile
					data, err := assets.Asset("data/smallterrain-blank.terrain")
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					} else {
						t.body = data
						goto send_tile
					}
				} else {
					http.Error(w, errors.New("The terrain tile does not exist").Error(), http.StatusNotFound)
				}
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// send the tile to the client
	send_tile:
		headers := w.Header()
		headers.Set("Content-Type", "application/octet-stream")
		headers.Set("Content-Encoding", "gzip")
		headers.Set("Content-Disposition", "attachment;filename="+vars["y"]+".terrain")
		w.Write(t.body)
	}
}

// An HTTP handler which returns a tileset's `layer.json` file
func layerHandler(tilesetRoot string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename := filepath.Join(tilesetRoot, vars["tileset"], "layer.json")

		// try and read the `layer.json` from disk
		body, err := ioutil.ReadFile(filename)
		if err != nil {
			if os.IsNotExist(err) {
				// check whether the tile directory exists
				_, err := os.Stat(filepath.Join(tilesetRoot, vars["tileset"]))
				if err != nil {
					if os.IsNotExist(err) {
						http.Error(w,
							fmt.Errorf("The tileset `%s` does not exist", vars["tileset"]).Error(),
							http.StatusNotFound)
						return
					} else {
						// There's some other problem (e.g. permissions)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}

				// the directory exists: send the default `layer.json`
				body = []byte(`{
  "tilejson": "2.1.0",
  "format": "heightmap-1.0",
  "version": "1.0.0",
  "scheme": "tms",
  "tiles": ["{z}/{x}/{y}.terrain"]
}`)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		headers := w.Header()
		headers.Set("Content-Type", "application/json")
		w.Write(body)
	}
}

// Return HTTP middleware which allows CORS requests from any domain
func addCorsHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := flag.Uint("port", 8000, "the port on which the server listens")
	tilesetRoot := flag.String("dir", ".", "the root directory under which tileset directories reside")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/tilesets/{tileset}/layer.json", layerHandler(*tilesetRoot))
	r.HandleFunc("/tilesets/{tileset}/{z:[0-9]+}/{x:[0-9]+}/{y:[0-9]+}.terrain", terrainHandler(*tilesetRoot))

	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, addCorsHeader(r)))

	log.Println("Terrain server listening on port", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
