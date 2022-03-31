package main

import (
	"fmt"

	"github.com/lispad/go-generics-tools/binheap"
)

func main() {
	someData := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	fmt.Printf("--- top 3 min ---\n source slice: %v\n", someData)
	mins := binheap.MinN[float64](someData, 3)
	fmt.Printf(" top 3 min elements: %v\n", mins)
	fmt.Printf(" slice was mutated: %v\n\n", someData)

	someIntData := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	fmt.Printf("--- top 4 max ---\n source slice: %v\n", someIntData)
	maxs := binheap.MaxN[int64](someIntData, 4)
	fmt.Printf(" top 4 max elements: %v\n\n", maxs)
	fmt.Printf(" slice was mutated: %v\n\n", someIntData)

	someFloat32Data := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	fmt.Printf("--- top 5 min without slice mutation ---\n source slice: %v\n", someFloat32Data)
	minsFloat32 := binheap.MinNImmutable[float32](someFloat32Data, 5)
	fmt.Printf(" top 5 min elements: %v\n", minsFloat32)
	fmt.Printf(" slice was not mutated: %v\n\n", someFloat32Data)

	someInt32Data := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	fmt.Printf("--- top 4 by remainder of devicion to 3 ---\n source slice: %v\n", someInt32Data)
	topByDivisionTo3Remainder := binheap.TopN[int64](someIntData, 4, func(x, y int64) bool { return x%3 > y%3 })
	fmt.Printf(" top 4 elements: %v\n", topByDivisionTo3Remainder)
	fmt.Printf(" slice was mutated: %v\n\n", someIntData)
}
