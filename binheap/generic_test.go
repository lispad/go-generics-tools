package binheap_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lispad/go-generics-tools/binheap"
)

func TestEmptyHeap(t *testing.T) {
	h := binheap.EmptyHeap[string](func(a, b string) bool { return len(a) > len(b) })
	assert.Equal(t, 0, h.Len())
	h.Push("1")
	assert.Equal(t, 1, h.Len())
	h.Push("22")
	assert.Equal(t, 2, h.Len())
	h.Push("4444")
	assert.Equal(t, 3, h.Len())
	h.Push("88888888")
	assert.Equal(t, 4, h.Len())
	h.Push("")
	assert.Equal(t, 5, h.Len())
	h.Push("55555")
	assert.Equal(t, 6, h.Len())
	h.Push("1")
	assert.Equal(t, 7, h.Len())
	h.Push("7777777")
	assert.Equal(t, 8, h.Len())
	h.Push("999999999")
	assert.Equal(t, 9, h.Len())
	h.Push("333")
	assert.Equal(t, 10, h.Len())

	assert.Equal(t, "999999999", h.Peak())
	assert.Equal(t, 10, h.Len())
	assert.Equal(t, "999999999", h.PushPop("4444")) // push less than max value, max will be returned
	assert.Equal(t, 10, h.Len())
	assert.Equal(t, "88888888", h.Peak())
	assert.Equal(t, 10, h.Len())
	assert.Equal(t, "0000000000", h.PushPop("0000000000")) // value will be returned, heap is unchanged
	assert.Equal(t, 10, h.Len())
	assert.Equal(t, "88888888", h.Peak())
	assert.Equal(t, 10, h.Len())

	assert.Equal(t, "88888888", h.Replace("22"))
	assert.Equal(t, 10, h.Len())
	assert.Equal(t, "7777777", h.Pop())
	assert.Equal(t, 9, h.Len())
	assert.Equal(t, "55555", h.Pop())
	assert.Equal(t, 8, h.Len())
	assert.Equal(t, "4444", h.Pop())
	assert.Equal(t, 7, h.Len())
	assert.Equal(t, "4444", h.Pop())
	assert.Equal(t, 6, h.Len())
	assert.Equal(t, "333", h.Pop())
	assert.Equal(t, 5, h.Len())
	h.Push("999999999")
	assert.Equal(t, 6, h.Len())
	assert.Equal(t, "999999999", h.Pop())
	assert.Equal(t, 5, h.Len())
	assert.Equal(t, "22", h.Pop())
	assert.Equal(t, 4, h.Len())
	assert.Equal(t, "22", h.Pop())
	assert.Equal(t, 3, h.Len())
	assert.Equal(t, "1", h.Pop())
	assert.Equal(t, 2, h.Len())
	assert.Equal(t, "1", h.Pop())
	assert.Equal(t, 1, h.Len())
	assert.Equal(t, "", h.Pop())
	assert.Equal(t, 0, h.Len())
}
