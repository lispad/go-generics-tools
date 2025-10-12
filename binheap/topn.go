package binheap

import (
	"cmp"
	"slices"
)

// TopN mutates data slice, moving TopN elements of slice to the beginning and returns subslice of first n elements.
// O(M *ln(N)), where M is slice size
// Mutates source data slice. No additional allocations are done.
func TopN[T any](data []T, n int, comparator CmpFunc[T]) []T {
	if n > len(data) {
		n = len(data)
	}
	heap := FromSlice(data[0:n], reverse(comparator))
	for i := n; i < len(data); i++ {
		if comparator(data[0], data[i]) == -1 {
			data[0], data[i] = data[i], data[0]
			heap.fix(0)
		}
	}
	slices.SortFunc(data[0:n], reverse(comparator))
	return data[0:n]
}

// MinN moves N Min elements of slice to the beginning and returns subslice of first n elements.
// Mutates source data slice.
// No additional allocations are done.
func MinN[T cmp.Ordered](data []T, n int) []T {
	return TopN(data, n, reverse(cmp.Compare[T]))
}

// MaxN moves N Max elements of slice to the beginning and returns subslice of first n elements.
// Mutates source data slice.
// No additional allocations are done.
func MaxN[T cmp.Ordered](data []T, n int) []T {
	return TopN(data, n, cmp.Compare[T])
}

func reverse[T any](comparator CmpFunc[T]) CmpFunc[T] {
	return func(x, y T) int {
		return -comparator(x, y)
	}
}
