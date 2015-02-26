package handlers

import (
	"errors"
	"fmt"
	"github.com/geo-data/cesium-terrain-server/assets"
	"github.com/geo-data/cesium-terrain-server/log"
	db "github.com/geo-data/cesium-terrain-server/stores"
	"github.com/geo-data/cesium-terrain-server/stores/items/terrain"
	"github.com/geo-data/cesium-terrain-server/stores/tiles"
	"github.com/gorilla/mux"
	"net/http"
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
