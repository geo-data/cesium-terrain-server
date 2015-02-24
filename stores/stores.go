package stores

import (
	"encoding"
	"errors"
)

var ErrNoItem = errors.New("item not found")

type Storer interface {
	Load(key string, obj encoding.BinaryUnmarshaler) error
	Save(key string, obj encoding.BinaryMarshaler) error
}
