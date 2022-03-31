package binheap_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"

	"github.com/lispad/go-generics-tools/binheap"
)

func TestMinHeapFromSlice(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		heapWithData := binheap.MinHeapFromSlice([]uint8{5, 10, 1, 5, 2, 7})
		top := make([]uint8, 0)
		for heapWithData.Len() > 0 {
			top = append(top, heapWithData.Pop())
		}
		assert.Equal(t, []uint8{1, 2, 5, 5, 7, 10}, top)
	})

	t.Run("strings", func(t *testing.T) {
		heapWithData := binheap.MinHeapFromSlice([]string{"foo", "bar", "foobar", "zzz", "aaa"})
		top := make([]string, 0)
		for heapWithData.Len() > 0 {
			top = append(top, heapWithData.Pop())
		}
		assert.Equal(t, []string{"aaa", "bar", "foo", "foobar", "zzz"}, top)
	})
}

func TestMaxHeapFromSlice(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		heapWithData := binheap.MaxHeapFromSlice([]uint8{5, 10, 1, 5, 2, 7})
		top := make([]uint8, 0)
		for heapWithData.Len() > 0 {
			top = append(top, heapWithData.Pop())
		}
		assert.Equal(t, []uint8{10, 7, 5, 5, 2, 1}, top)
	})

	t.Run("strings", func(t *testing.T) {
		heapWithData := binheap.MaxHeapFromSlice([]string{"foo", "bar", "foobar", "zzz", "aaa"})
		top := make([]string, 0)
		for heapWithData.Len() > 0 {
			top = append(top, heapWithData.Pop())
		}
		assert.Equal(t, []string{"zzz", "foobar", "foo", "bar", "aaa"}, top)
	})
}

func FuzzEmptyMinHeap(f *testing.F) {
	min := binheap.EmptyMinHeap[byte]()
	f.Fuzz(func(t *testing.T, data []byte) {
		for _, b := range data {
			min.Push(b)
		}

		sorted := make([]byte, len(data))
		copy(sorted, data)
		slices.Sort(sorted)
		for _, b := range sorted {
			v := min.Pop()
			assert.Equal(t, b, v)
		}
		assert.Equal(t, 0, min.Len())
	})
}

func FuzzEmptyMaxHeap(f *testing.F) {
	max := binheap.EmptyMaxHeap[byte]()
	f.Fuzz(func(t *testing.T, data []byte) {
		for _, b := range data {
			max.Push(b)
		}
		reverseSorted := make([]byte, len(data))
		copy(reverseSorted, data)
		sort.Slice(reverseSorted, func(i, j int) bool {
			return reverseSorted[i] > reverseSorted[j]
		})
		for _, b := range reverseSorted {
			v := max.Pop()
			assert.Equal(t, b, v)
		}
		assert.Equal(t, 0, max.Len())
	})
}

func TestEmptyComparableHeapHeap(t *testing.T) {
	h := binheap.EmptyComparableHeap[string](func(a, b string) bool { return len(a) > len(b) })
	assert.Equal(t, 0, h.Len())
	assert.False(t, h.Search("333"))
	assert.False(t, h.Delete("333"))

	h.Push("1")
	assert.Equal(t, 1, h.Len())
	assert.False(t, h.Search("333"))
	assert.False(t, h.Delete("333"))
	assert.True(t, h.Search("1"))

	h.Push("22")
	assert.Equal(t, 2, h.Len())
	assert.True(t, h.Search("1"))
	assert.True(t, h.Delete("1"))
	assert.False(t, h.Search("1"))
	assert.True(t, h.Search("22"))
	assert.Equal(t, 1, h.Len())

	h.Push("4444")
	h.Push("88888888")
	h.Push("55555")
	h.Push("1")
	h.Push("4444")
	h.Push("4444")
	h.Push("7777777")
	assert.Equal(t, 8, h.Len())

	assert.True(t, h.Search("88888888"))
	assert.True(t, h.Delete("88888888")) // test root deletion
	assert.False(t, h.Search("88888888"))
	assert.False(t, h.Delete("88888888"))

	assert.True(t, h.Search("4444")) // first entry
	assert.True(t, h.Delete("4444"))
	assert.True(t, h.Search("4444")) // second entry
	assert.True(t, h.Delete("4444"))
	assert.True(t, h.Search("4444")) // third entry
	assert.True(t, h.Delete("4444"))
	assert.False(t, h.Search("4444")) // no more "4444"
	assert.False(t, h.Delete("4444"))
}
