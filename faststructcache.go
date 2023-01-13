package faststructcache

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/kelindar/binary"
	"github.com/klauspost/compress/s2"
)

type (
	// FastStructCache is an in-memory key-value store intended to be able to store larger
	// structs (although any data type should work).
	// It wraps fastcache from VictoriaMetrics with the binary codec from kelindar,
	// and optional snappy-like compression (S2 from klauspost).
	FastStructCache[KT, VT any] struct {
		Cache    *fastcache.Cache
		compress bool
	}
)

// New returns new cache with the given maxBytes capacity in bytes (min 32MB), with optional light compression.
func New[KT, VT any](mb int, compress bool) FastStructCache[KT, VT] {
	return FastStructCache[KT, VT]{Cache: fastcache.New(mb * 1024 * 1024), compress: compress}
}

// Set stores (k, v) in the cache.
func (c FastStructCache[KT, VT]) Set(k KT, v VT) {
	c.Cache.Set(c.encode(k, v))
}

// Get returns pointer to a value stored for given key k, or nil if it didn't exist in the cache
func (c FastStructCache[KT, VT]) Get(k KT) *VT {
	encodedVal, exists := c.Cache.HasGet(nil, c.encodeKey(k))
	if !exists {
		return nil
	}
	return c.decodeVal(encodedVal)
}

// Has returns true if entry for the given key k exists in the cache.
func (c FastStructCache[KT, VT]) Has(k KT) bool {
	return c.Cache.Has(c.encodeKey(k))
}

// Del deletes value for the given k from the cache.
func (c FastStructCache[KT, VT]) Del(k KT) {
	c.Cache.Del(c.encodeKey(k))
}

// Reset removes all the items from the cache.
func (c FastStructCache[KT, VT]) Reset() {
	c.Cache.Reset()
}

func (c FastStructCache[KT, VT]) encode(k KT, v VT) ([]byte, []byte) {
	encodedVal, err := binary.Marshal(v)
	if c.compress {
		encodedVal = s2.Encode(nil, encodedVal)
	}
	if err != nil {
		panic(err)
	}
	return c.encodeKey(k), encodedVal
}

func (c FastStructCache[KT, VT]) encodeKey(k KT) []byte {
	encodedKey, err := binary.Marshal(k)
	if err != nil {
		panic(err)
	}
	return encodedKey
}

func (c FastStructCache[KT, VT]) decodeVal(encodedVal []byte) *VT {
	if c.compress {
		var err error
		encodedVal, err = s2.Decode(nil, encodedVal)
		if err != nil {
			panic(err)
		}
	}

	var v VT
	err := binary.Unmarshal(encodedVal, &v)
	if err != nil {
		panic(err)
	}
	return &v
}
