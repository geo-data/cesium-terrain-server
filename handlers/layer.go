package handlers

import (
	"fmt"
	"net/http"

	"github.com/geo-data/cesium-terrain-server/log"
	"github.com/geo-data/cesium-terrain-server/stores"
	"github.com/gorilla/mux"
)

// An HTTP handler which returns a tileset's `layer.json` file
func LayerHandler(store stores.Storer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err   error
			layer []byte
		)

		defer func() {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Err(err.Error())
			}
		}()

		vars := mux.Vars(r)

		// Try and get a `layer.json` from the stores
		layer, err = store.Layer(vars["tileset"])
		if err == stores.ErrNoItem {
			err = nil // don't persist this error
			if store.TilesetStatus(vars["tileset"]) == stores.NOT_FOUND {
				http.Error(w,
					fmt.Errorf("The tileset `%s` does not exist", vars["tileset"]).Error(),
					http.StatusNotFound)
				return
			}

			// the directory exists: send the default `layer.json`
			layer = []byte(`{
  "tilejson": "2.1.0",
  "format": "heightmap-1.0",
  "version": "1.0.0",
  "scheme": "tms",
  "tiles": ["{z}/{x}/{y}.terrain"]
}`)
		} else if err != nil {
			return
		}

		headers := w.Header()
		headers.Set("Content-Type", "application/json")
		w.Write(layer)
	}
}
