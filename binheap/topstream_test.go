package binheap_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"

	"github.com/lispad/go-generics-tools/binheap"
)

func TestMaxNImmutable(t *testing.T) {
	source := []string{"foo", "bar", "baz", "zzz", "aaa"}
	input := make([]string, len(source))
	copy(input, source)

	max := binheap.MaxNImmutable(input, 3)
	assert.Equal(t, source, input) // input should not be mutated
	assert.Equal(t, []string{"zzz", "foo", "baz"}, max)
}

func TestMinNImmutable(t *testing.T) {
	source := []string{"foo", "bar", "baz", "zzz", "aaa"}
	input := make([]string, len(source))
	copy(input, source)

	max := binheap.MinNImmutable(input, 3)
	assert.Equal(t, source, input) // input should not be mutated
	assert.Equal(t, []string{"aaa", "bar", "baz"}, max)
}

func FuzzTopNImmutable(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte, i uint) {
		n := int(i)
		if n == 0 {
			n = 1
		}
		if n > len(data) {
			n = len(data)
		}

		input := make([]byte, len(data))
		copy(input, data)

		sorted := make([]byte, len(data))
		copy(sorted, data)
		slices.Sort(sorted)

		topN := binheap.MinNImmutable(input, n)
		assert.Equal(t, sorted[0:n], topN) // top should match with order
		assert.Equal(t, data, input)       // input should not be mutated

		topN = binheap.MaxNImmutable(input, n)
		// reverse sort
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i] > sorted[j]
		})
		assert.Equal(t, sorted[0:n], topN) // top should match with order
		assert.Equal(t, data, input)       // input should not be mutated
	})
}
