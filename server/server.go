package server

import (
	"errors"
	"fmt"
	"github.com/geo-data/cesium-terrain-server/assets"
	"github.com/geo-data/cesium-terrain-server/log"
	db "github.com/geo-data/cesium-terrain-server/stores"
	"github.com/geo-data/cesium-terrain-server/stores/items"
	"github.com/geo-data/cesium-terrain-server/stores/items/terrain"
	"github.com/geo-data/cesium-terrain-server/stores/tiles"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
)

// An HTTP handler which returns a terrain tile resource
func TerrainHandler(stores []*tiles.Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			t   terrain.Terrain
			err error
		)

		defer func() {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Err(err.Error())
			}
		}()

		// get the tile coordinate from the URL
		vars := mux.Vars(r)
		err = t.ParseCoord(vars["x"], vars["y"], vars["z"])
		if err != nil {
			return
		}

		// Try and get a tile from the stores
		var idx int
		for i, store := range stores {
			idx = i
			err = store.LoadTile(vars["tileset"], &t)
			if err == nil {
				break
			} else if err == db.ErrNoItem {
				continue
			} else {
				return
			}
		}

		if err == db.ErrNoItem {
			if t.IsRoot() {
				// serve up a blank tile as it is a missing root tile
				data, err := assets.Asset("data/smallterrain-blank.terrain")
				if err != nil {
					return
				} else {
					err = t.UnmarshalBinary(data)
					if err != nil {
						return
					}
				}
			} else {
				err = nil
				http.Error(w, errors.New("The terrain tile does not exist").Error(), http.StatusNotFound)
				return
			}
		}

		body, err := t.MarshalBinary()
		if err != nil {
			return
		}

		// send the tile to the client
		headers := w.Header()
		headers.Set("Content-Type", "application/octet-stream")
		headers.Set("Content-Encoding", "gzip")
		headers.Set("Content-Disposition", "attachment;filename="+vars["y"]+".terrain")
		w.Write(body)

		// Save the tile in any preceding stores that didn't have it.
		if idx > 0 {
			for j := 0; j < idx; j++ {
				if err := stores[j].SaveTile(vars["tileset"], &t); err != nil {
					log.Err(fmt.Sprintf("failed to store tileset: %s", err))
				}
			}
		}
	}
}

// An HTTP handler which returns a tileset's `layer.json` file
func LayerHandler(tilesetRoot string, stores []db.Storer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err      error
			response items.Item
		)

		defer func() {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Err(err.Error())
			}
		}()

		vars := mux.Vars(r)
		key := filepath.Join(vars["tileset"], "layer.json")

		// Try and get a `layer.json` from the stores
		var idx int
		for i, store := range stores {
			idx = i
			err = store.Load(key, &response)
			if err == nil {
				break
			} else if err == db.ErrNoItem {
				continue
			} else {
				return
			}
		}

		if err == db.ErrNoItem {
			// check whether the tile directory exists
			_, err = os.Stat(filepath.Join(tilesetRoot, vars["tileset"]))
			if err != nil {
				if os.IsNotExist(err) {
					err = nil
					http.Error(w,
						fmt.Errorf("The tileset `%s` does not exist", vars["tileset"]).Error(),
						http.StatusNotFound)
					return
				}
				// There's some other problem (e.g. permissions)
				return
			}

			// the directory exists: send the default `layer.json`

			err = response.UnmarshalBinary([]byte(`{
  "tilejson": "2.1.0",
  "format": "heightmap-1.0",
  "version": "1.0.0",
  "scheme": "tms",
  "tiles": ["{z}/{x}/{y}.terrain"]
}`))
			if err != nil {
				return
			}
		}

		body, err := response.MarshalBinary()
		if err != nil {
			return
		}

		headers := w.Header()
		headers.Set("Content-Type", "application/json")
		w.Write(body)

		// Save the json file in any preceding stores that didn't have it.
		if idx > 0 {
			for j := 0; j < idx; j++ {
				if err := stores[j].Save(key, &response); err != nil {
					log.Err(fmt.Sprintf("failed to store layer.json: %s", err))
				}
			}
		}
	}
}

// Return HTTP middleware which allows CORS requests from any domain
func AddCorsHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
