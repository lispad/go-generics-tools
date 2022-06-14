# GoLang Generics tools: Heap structure, sharded rw-locked map.
[![Go Report Card](https://goreportcard.com/badge/github.com/lispad/go-generics-tools)](https://goreportcard.com/report/github.com/lispad/go-generics-tools)
[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Introduction
------------

The [Heap](binheap/README.md) package contains simple [binary heap](https://en.wikipedia.org/wiki/Binary_heap) implementation, using Golang
generics. There are several heap implementations
[Details](binheap/README.md).


The [ShardedLockMap](smap/README.md) package contains implementation of sharded lock map.
Interface is similar to sync.map, but sharded lock map is faster on scenarios with huge read load with rare updates,
and uses less memory, doing less allocations.
[Details](smap/README.md)

Compatibility
-------------
Minimal Golang version is 1.18. Generics and fuzz testing are used.

Installation
----------------------

To install package, run:

    go get github.com/lispad/go-generics-tools/binheap
or

    go get github.com/lispad/go-generics-tools/smap


License
-------

The binheap package is licensed under the MIT license. Please see the LICENSE file for details.