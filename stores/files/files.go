package files

import (
	"encoding"
	"github.com/geo-data/cesium-terrain-server/stores"
	"io/ioutil"
	"log"
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
	log.Printf("save fs: %s", key)
	return nil
}

// Load a terrain tile on disk into the Terrain structure.
func (this *Store) Load(key string, obj encoding.BinaryUnmarshaler) (err error) {
	log.Printf("load fs key: %s", key)
	filename := filepath.Join(
		this.root,
		key)
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			err = stores.ErrNoItem
		}
		return
	}

	err = obj.UnmarshalBinary(body)
	log.Printf("load fs: %s", filename)
	return
}
