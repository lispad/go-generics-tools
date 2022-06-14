# Sharded RWLocked Map, using go generics
[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Introduction
------------

In some scenarios sync.Map could allocate large space for dirty-read copies, or even heavily use locks (when getting 
missing values, that should be rechecked in map with lock). In such scenarios go internal map with rw-mutex, divided to 
several shards could perform much better, with cost of memory for additional locks storage. But this amount could be
much less, than sync.map uses.

Interface incompatibility
------------

- `LoadOrStore` method changed to `LoadOrCreate`, with callback that generates value. Could be used to avoid 
unnecessary creating huge values, in case if key already exists. 

Usage Example
-----------------

    package main
    
    import (
        "fmt"

        "github.com/lispad/go-generics-tools/smap"
    )
    
    func main() {
    	m := smap.NewIntegerComparable[int, int](8, 128)
        m.Store(123, 456)

        value, ok := m.Load(123)
        fmt.Printf("%d, %t", value, ok)
    }

A bit more examples could be found in tests.

Benchmark
-----------------


#### Benchmark

Benchmark performed on Lenovo Ideapad laptop with AMD Ryzen 7 4700U, Linux Mint 20.3 with 5.13.0 kernel

    BenchmarkIntegerSMap_ConcurrentGet-8              	82540485	     12.95 ns/op	    0 B/op	    0 allocs/op
    BenchmarkSyncMap_ConcurrentGet-8                  	73431339	     18.35 ns/op	    0 B/op	    0 allocs/op
    BenchmarkLockMap_ConcurrentGet-8                  	19327282	     54.03 ns/op	    0 B/op	    0 allocs/op
    BenchmarkIntegerShardedMap_ConcurrentSet-8        	25605380	     42.33 ns/op	    0 B/op	    0 allocs/op
    BenchmarkSyncMap_ConcurrentSet-8                  	 2138496	     536.1 ns/op	   36 B/op	    3 allocs/op
    BenchmarkLockMap_ConcurrentSet-8                  	 3827476	     302.9 ns/op	    0 B/op	    0 allocs/op
    BenchmarkIntegerShardedMap_ConcurrentGetSet5-8    	 3377473	     357.1 ns/op	    0 B/op	    0 allocs/op
    BenchmarkSyncMap_ConcurrentGetSet5-8              	  323318	      3540 ns/op	  266 B/op	    3 allocs/op
    BenchmarkLockMap_ConcurrentGetSet5-8              	  204633	      6176 ns/op	    0 B/op	    0 allocs/op
    BenchmarkIntegerShardedMap_ConcurrentGetSet50-8   	18236305	     65.03 ns/op	    0 B/op	    0 allocs/op
    BenchmarkSyncMap_ConcurrentGetSet50-8             	18337423	     56.55 ns/op	   32 B/op	    2 allocs/op
    BenchmarkLockMap_ConcurrentGetSet50-8             	 2697315	     431.2 ns/op	    0 B/op	    0 allocs/op
    BenchmarkIntegerShardedMap_ConcurrentGetSet1-8    	 1003506	      1076 ns/op	    0 B/op	    0 allocs/op
    BenchmarkSyncMap_ConcurrentGetSet1-8              	   32668	     44423 ns/op	 3508 B/op	    5 allocs/op
    BenchmarkLockMap_ConcurrentGetSet1-8              	   78486	     16019 ns/op	    0 B/op	    0 allocs/op

Sharded Lock map is approximately equal to sync.Map on 50% read + 50% concurrent writes, and is much faster on 
5% writes+95% reads, and 1% writes+99% reads. 
Also sharded map allocated about 40x less memory in 5% writes+95% reads scenario, than sync.Map does.

Compatibility
-------------
Minimal Golang version is 1.18. Generics are used.

Installation
----------------------

To install package, run:

    go get github.com/lispad/go-generics-tools/smap

License
-------

The smap package is licensed under the MIT license. Please see the LICENSE file for details.
