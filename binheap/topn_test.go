package binheap_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"

	"github.com/lispad/go-generics-tools/binheap"
)

func TestMinN(t *testing.T) {
	t.Run("float32", func(t *testing.T) {
		source := []float32{11.3, 5, 10.12, 1, 5, 2, 7}
		input := make([]float32, len(source))
		copy(input, source)

		min3 := binheap.MinN(input, 3)
		assert.Equal(t, []float32{1, 2, 5}, min3)
		assert.Equal(t, []float32{1, 2, 5, 11.3, 10.12, 5, 7}, input) // check that rest of elements are kept

		min5 := binheap.MinN(input, 5)
		assert.Equal(t, []float32{1, 2, 5, 5, 7}, min5)
		assert.ElementsMatch(t, input, source)

		min10 := binheap.MinN(input, 10)
		assert.Equal(t, []float32{1, 2, 5, 5, 7, 10.12, 11.3}, min10)
		assert.ElementsMatch(t, input, source)
	})

	t.Run("strings", func(t *testing.T) {
		source := []string{"foo", "bar", "foobar", "zzz", "aaa", "some more text"}
		input := make([]string, len(source))
		copy(input, source)

		min1 := binheap.MinN(input, 1)
		assert.Equal(t, []string{"aaa"}, min1)
		assert.ElementsMatch(t, input, source)

		min3 := binheap.MinN(input, 3)
		assert.Equal(t, []string{"aaa", "bar", "foo"}, min3)
		assert.ElementsMatch(t, input, source)
	})
}

func TestMaxN(t *testing.T) {
	t.Run("int16", func(t *testing.T) {
		source := []int16{5, 1, 5, 2, 7, -111}
		input := make([]int16, len(source))
		copy(input, source)

		max3 := binheap.MaxN(input, 3)
		assert.Equal(t, []int16{7, 5, 5}, max3)
		assert.ElementsMatch(t, input, source)

		max5 := binheap.MaxN(input, 5)
		assert.Equal(t, []int16{7, 5, 5, 2, 1}, max5)
		assert.ElementsMatch(t, input, source)

		max10 := binheap.MaxN(input, 10)
		assert.Equal(t, []int16{7, 5, 5, 2, 1, -111}, max10)
		assert.ElementsMatch(t, input, source)
	})

	t.Run("strings", func(t *testing.T) {
		source := []string{"foo", "bar", "foobar", "zzz", "aaa", "some more text"}
		input := make([]string, len(source))
		copy(input, source)

		max1 := binheap.MaxN(input, 1)
		assert.Equal(t, []string{"zzz"}, max1)
		assert.ElementsMatch(t, input, source)

		max3 := binheap.MaxN(input, 3)
		assert.Equal(t, []string{"zzz", "some more text", "foobar"}, max3)
		assert.ElementsMatch(t, input, source)
	})
}

func FuzzTopN(f *testing.F) {
	f.Add([]byte{1, 10, 21, 211, 12, 2, 7, 13, 10}, uint(5))
	f.Fuzz(func(t *testing.T, data []byte, i uint) {
		n := int(i)
		if n > len(data) {
			n = len(data)
		}

		sorted := make([]byte, len(data))
		copy(sorted, data)
		slices.Sort(sorted)

		topN := binheap.MinN(data, n)
		assert.Equal(t, sorted[0:n], topN)            // top should match with order
		assert.ElementsMatch(t, sorted[n:], data[n:]) // rest of elements should match but could be unsorted

		slices.Sort(data)
		topN = binheap.MinN(data, n)
		assert.Equal(t, sorted[0:n], topN) // top should match with order
		assert.Equal(t, sorted, data)      // data should not be changed if it was already sorted
	})
}
