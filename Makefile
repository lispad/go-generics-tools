test: fuzz
	go test -cover ./...

fuzz:
	go test -fuzz=FuzzEmptyMinHeap -fuzztime=10s github.com/lispad/go-generics-tools/binheap
	go test -fuzz=FuzzEmptyMaxHeap -fuzztime=10s github.com/lispad/go-generics-tools/binheap
	go test -fuzz=FuzzTopN$$ -fuzztime=10s github.com/lispad/go-generics-tools/binheap
	go test -fuzz=FuzzTopNImmutable -fuzztime=10s github.com/lispad/go-generics-tools/binheap

lint:
	golangci-lint run