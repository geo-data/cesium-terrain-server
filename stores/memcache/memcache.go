package memcache

import (
	"encoding"
	mc "github.com/bradfitz/gomemcache/memcache"
	"github.com/geo-data/cesium-terrain-server/stores"
	"log"
)

type Store struct {
	mc *mc.Client
}

func New(connstr string) stores.Storer {
	return &Store{
		mc: mc.New(connstr),
	}
}

func (this *Store) Save(key string, obj encoding.BinaryMarshaler) (err error) {
	log.Printf("save mem: %s", key)
	value, err := obj.MarshalBinary()
	if err != nil {
		return
	}
	return this.mc.Set(&mc.Item{Key: key, Value: value})
}

func (this *Store) Load(key string, obj encoding.BinaryUnmarshaler) (err error) {
	val, err := this.mc.Get(key)
	if err != nil {
		if err == mc.ErrCacheMiss {
			log.Printf("load mem err: %s", err)
			err = stores.ErrNoItem
		}
		return
	}
	log.Printf("load mem: %s", key)
	err = obj.UnmarshalBinary(val.Value)
	return
}
