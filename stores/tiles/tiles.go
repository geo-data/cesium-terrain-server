package tiles

import (
	"encoding"
	"github.com/geo-data/cesium-terrain-server/stores"
	"github.com/geo-data/cesium-terrain-server/stores/items/terrain"
)

type Namer interface {
	TileName(tileset string, tile *terrain.Terrain) string
}

type Store struct {
	Namer Namer
	Store stores.Storer
}

func New(namer Namer, store stores.Storer) *Store {
	return &Store{
		Namer: namer,
		Store: store,
	}
}

func (this *Store) LoadTile(tileset string, tile *terrain.Terrain) error {
	key := this.Namer.TileName(tileset, tile)
	return this.Store.Load(key, tile)
}

func (this *Store) SaveTile(tileset string, tile *terrain.Terrain) error {
	key := this.Namer.TileName(tileset, tile)
	return this.Store.Save(key, tile)
}

func (this *Store) Save(key string, obj encoding.BinaryMarshaler) error {
	return this.Store.Save(key, obj)
}

func (this *Store) Load(key string, obj encoding.BinaryUnmarshaler) error {
	return this.Store.Load(key, obj)
}
