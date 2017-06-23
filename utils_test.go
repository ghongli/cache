package cache

import (
	"reflect"
	"testing"
	"time"
)

func TestBytesToItem(t *testing.T) {
	serializer := NewCacheSerializer()

	item := CacheItem{
		CreatedTime: time.Now(),
		Expired:     time.Now().Sub(time.Now()),
		Data:        []byte("pong"),
	}
	b, err := serializer.Serialized(item)
	if err != nil {
		t.Error(err)
	}

	i := new(CacheItem)
	err = serializer.DeSerialized(b, i)
	if err != nil {
		t.Error(err)
	}

	if !i.CreatedTime.Add(i.Expired).Equal(item.CreatedTime.Add(item.Expired)) {
		t.Fatalf("Time should equal.. Expected %v \n Got %v", item.CreatedTime.Add(item.Expired), i.CreatedTime.Add(i.Expired))
	}

	if !reflect.DeepEqual(item.Data, i.Data) {
		t.Fatalf("Data not equal.. Expected %v \n. Got %v", item.Data, i.Data)
	}
}

func TestSerializer_SerializedIfTypeIsByteArray(t *testing.T) {
	serializer := NewCacheSerializer()

	byt := make([]byte, 5)
	newByt, err := serializer.Serialized(byt)
	if err != nil {
		t.Fatalf("Expected serialization to be successful... %v", err)
	}

	if !reflect.DeepEqual(byt, newByt) {
		t.Fatalf("Values differ \n Expected %v \n Got %v", byt, newByt)
	}
}
