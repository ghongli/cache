package cache

import (
	"bytes"
	"encoding/gob"
)

type Serializer interface {
	Serialized(i interface{}) ([]byte, error)
	DeSerialized(data []byte, i interface{}) error
}

func NewCacheSerializer() Serializer {
	return &SerializeHelper{}
}

// a helper types to serialize and deserialize
type SerializeHelper struct {
}

// Convert a given type into a byte array
// use the encoding/gob package
func (sh *SerializeHelper) Serialized(i interface{}) ([]byte, error) {
	if sh, ok := i.([]byte); ok {
		return sh, nil
	}

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	if err := enc.Encode(i); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Writes a byte array into a type.
// use the encoding/gob package
func (sh *SerializeHelper) DeSerialized(data []byte, i interface{}) error {
	if sh, ok := i.(*[]byte); ok {
		*sh = data
		return nil
	}

	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(i)
}
