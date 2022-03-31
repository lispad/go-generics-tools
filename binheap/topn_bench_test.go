package binheap_test

import (
	"math/rand"
	"testing"

	"golang.org/x/exp/slices"

	"github.com/lispad/go-generics-tools/binheap"
)

func BenchmarkSortedMaxN(b *testing.B) {
	sl := testSlice(b.N)
	k := 5
	b.ReportAllocs()
	b.ResetTimer()

	if k > len(sl) {
		k = len(sl)
	}
	slices.Sort(sl)
	result := sl[:k]
	_ = result
}

func BenchmarkMaxNImmutable(b *testing.B) {
	sl := testSlice(b.N)
	k := 5
	b.ReportAllocs()
	b.ResetTimer()

	if k > len(sl) {
		k = len(sl)
	}
	result := binheap.MaxNImmutable(sl, k)
	_ = result
}

func BenchmarkMaxN(b *testing.B) {
	sl := testSlice(b.N)
	k := 5
	b.ReportAllocs()
	b.ResetTimer()

	if k > len(sl) {
		k = len(sl)
	}
	result := binheap.MinN(sl, k)
	_ = result
}

func testSlice(size int) []int32 {
	sl := make([]int32, size)
	for i := 0; i < size; i++ {
		sl[i] = rand.Int31()
	}
	return sl
}
