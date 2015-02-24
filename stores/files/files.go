package files

import (
	"encoding"
	"fmt"
	"github.com/geo-data/cesium-terrain-server/log"
	"github.com/geo-data/cesium-terrain-server/stores"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Store struct {
	root string
}

func New(root string) stores.Storer {
	return &Store{
		root: root,
	}
}

// This is a no-op
func (this *Store) Save(key string, obj encoding.BinaryMarshaler) error {
	log.Debug(fmt.Sprintf("file store: save: %s", key))
	return nil
}

// Load a terrain tile on disk into the Terrain structure.
func (this *Store) Load(key string, obj encoding.BinaryUnmarshaler) (err error) {
	filename := filepath.Join(
		this.root,
		key)

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug(fmt.Sprintf("file store: not found: %s", filename))
			err = stores.ErrNoItem
		}
		return
	}

	log.Debug(fmt.Sprintf("file store: load: %s", filename))
	err = obj.UnmarshalBinary(body)
	return
}
