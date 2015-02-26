package fs

import (
	"fmt"
	"github.com/geo-data/cesium-terrain-server/log"
	"github.com/geo-data/cesium-terrain-server/stores"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type Store struct {
	root string
}

func New(root string) stores.Storer {
	return &Store{
		root: root,
	}
}

func (this *Store) readFile(filename string) (body []byte, err error) {
	body, err = ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug(fmt.Sprintf("file store: not found: %s", filename))
			err = stores.ErrNoItem
		}
		return
	}

	log.Debug(fmt.Sprintf("file store: load: %s", filename))
	return
}

// Load a terrain tile on disk into the Terrain structure.
func (this *Store) Tile(tileset string, tile *stores.Terrain) (err error) {
	filename := filepath.Join(
		this.root,
		tileset,
		strconv.FormatUint(tile.Z, 10),
		strconv.FormatUint(tile.X, 10),
		strconv.FormatUint(tile.Y, 10)+".terrain")

	body, err := this.readFile(filename)
	if err != nil {
		return
	}

	err = tile.UnmarshalBinary(body)
	return
}

func (this *Store) Layer(tileset string) ([]byte, error) {
	filename := filepath.Join(tileset, "layer.json")
	return this.readFile(filename)
}

func (this *Store) TilesetStatus(tileset string) (status stores.TilesetStatus) {
	// check whether the tile directory exists
	_, err := os.Stat(filepath.Join(this.root, tileset))
	if err != nil {
		if os.IsNotExist(err) {
			return stores.FOUND
		}
	}

	return stores.NOT_FOUND
}
