package handlers

import (
	"errors"
	"fmt"
	"github.com/geo-data/cesium-terrain-server/assets"
	"github.com/geo-data/cesium-terrain-server/log"
	"github.com/geo-data/cesium-terrain-server/stores"
	"github.com/gorilla/mux"
	"net/http"
)

// An HTTP handler which returns a terrain tile resource
func TerrainHandler(store stores.Storer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			t   stores.Terrain
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

		// Try and get a tile from the store
		err = store.Tile(vars["tileset"], &t)
		if err == stores.ErrNoItem {
			if store.TilesetStatus(vars["tileset"]) == stores.NOT_FOUND {
				err = nil
				http.Error(w,
					fmt.Errorf("The tileset `%s` does not exist", vars["tileset"]).Error(),
					http.StatusNotFound)
				return
			}

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
		} else if err != nil {
			return
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
	}
}
