package main

import (
	"fmt"

	"github.com/lispad/go-generics-tools/binheap"
)

type someStruct struct {
	someData []int // make struct non comparable
}

func main() {
	heap := binheap.EmptyHeap[someStruct](func(x, y someStruct) bool { return len(x.someData) > len(y.someData) })
	heap.Push(someStruct{someData: []int{1, 2, 3, 4, 5, 6}})
	heap.Push(someStruct{someData: nil})
	heap.Push(someStruct{someData: []int{1, 2, 3}})
	heap.Push(someStruct{someData: []int{11}})
	heap.Push(someStruct{someData: []int{11, 2222, 3333, 44}})
	heap.Push(someStruct{someData: []int{1}})
	heap.Push(someStruct{someData: nil})
	heap.Push(someStruct{someData: []int{1111}})
	heap.Push(someStruct{someData: []int{1, 2}})

	fmt.Printf("Heap has len: %d\nElements sorted by slice length:\n\n", heap.Len())
	for heap.Len() > 0 {
		fmt.Printf(" %v\n", heap.Pop())
	}
}
