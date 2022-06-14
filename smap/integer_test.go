package smap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneric_Delete(t *testing.T) {
	m := NewInteger[uint32, string](8, 128)
	val, ok := m.Load(123)
	assert.False(t, ok)
	assert.Equal(t, "", val)

	m.Delete(123)         // no error on deleting
	val, ok = m.Load(123) // nothing changed
	assert.False(t, ok)
	assert.Equal(t, "", val)

	m.Store(123, "value set")
	val, ok = m.Load(123)
	assert.True(t, ok)
	assert.Equal(t, "value set", val)

	m.Delete(123)
	val, ok = m.Load(123)
	assert.False(t, ok)
	assert.Equal(t, "", val)
}

func TestGeneric_GetAndDelete(t *testing.T) {
	m := NewInteger[uint32, string](8, 128)
	m.Store(123, "value set")
	val, ok := m.Load(123)
	assert.True(t, ok)
	assert.Equal(t, "value set", val)

	val, ok = m.Load(123) // nothing changed
	assert.True(t, ok)
	assert.Equal(t, "value set", val)

	val, ok = m.LoadAndDelete(123)
	assert.True(t, ok)
	assert.Equal(t, "value set", val)

	val, ok = m.LoadAndDelete(123)
	assert.False(t, ok)
	assert.Equal(t, "", val)
}

func TestGeneric_GetOrCreate(t *testing.T) {
	m := NewInteger[uint32, string](8, 128)
	val, ok := m.Load(123)
	assert.False(t, ok)
	assert.Equal(t, "", val)

	val, ok = m.LoadOrCreate(123, func() string { return "new value created" })
	assert.False(t, ok)
	assert.Equal(t, "new value created", val)

	val, ok = m.LoadOrCreate(123, func() string { panic("no generator should be called") })
	assert.True(t, ok)
	assert.Equal(t, "new value created", val)

	val, ok = m.LoadAndDelete(123)
	assert.True(t, ok)
	assert.Equal(t, "new value created", val)
}

func TestGeneric_Range(t *testing.T) {
	m := NewInteger[int, int](8, 128)
	expected := make(map[int]int, 256)
	for i := 0; i < 256; i++ {
		m.Store(i, i*i)
		expected[i] = i * i
	}

	result := make(map[int]int, 256)
	m.Range(func(k int, v int) bool {
		result[k] = v
		return true
	})
	assert.Equal(t, expected, result)
}

func TestGenericComparable_CompareAndSwap(t *testing.T) {
	m := NewIntegerComparable[int, int](8, 128)
	m.CompareAndSwap(123, 0, 23) // no value with key 123 => no change
	val, ok := m.Load(123)
	assert.False(t, ok)
	assert.Equal(t, 0, val)

	m.Store(123, 10)
	m.CompareAndSwap(123, 11, 23) // value differs from 11 => no change
	val, ok = m.Load(123)
	assert.True(t, ok)
	assert.Equal(t, 10, val)

	m.CompareAndSwap(123, 10, 23) // value equal 11 => change to 23
	val, ok = m.Load(123)
	assert.True(t, ok)
	assert.Equal(t, 23, val)
}
