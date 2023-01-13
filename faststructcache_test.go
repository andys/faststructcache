package faststructcache

import (
	"testing"

	"github.com/zeebo/assert"
)

type TestKey struct {
	Key string
}

type TestVal struct {
	Number    uint32
	CopyOfKey TestKey
	Value     string
}

func TestCache(t *testing.T) {
	c := New[TestKey, TestVal](32, false)
	assert.NotNil(t, c)
	testCache(c, t)
}

func TestCacheCompressed(t *testing.T) {
	c := New[TestKey, TestVal](32, true)

	assert.NotNil(t, c)
	testCache(c, t)
}

func testCache(c FastStructCache[TestKey, TestVal], t *testing.T) {
	key := TestKey{Key: "MyKeyString123"}
	var nilVal *TestVal
	val := TestVal{Number: 456, CopyOfKey: key, Value: "PMyValueString789"}

	assert.Equal(t, c.Has(key), false)
	assert.Equal(t, c.Get(key), nilVal)
	c.Set(key, val)
	assert.Equal(t, c.Has(key), true)

	gotten := c.Get(key)
	assert.Equal(t, *gotten, val)

	c.Del(key)
	assert.Equal(t, c.Get(key), nilVal)
	assert.Equal(t, c.Has(key), false)
}
