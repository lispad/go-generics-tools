package smap

import (
	"golang.org/x/exp/constraints"
)

// NewInteger creates sharded rwlock maps with shard detection based on key division to shards count modulo.
func NewInteger[K constraints.Integer, V any](shardsCount, defaultSize int) Generic[K, V] {
	return NewGeneric[K, V](shardsCount, defaultSize, func(key K) int {
		return int(key) % shardsCount
	})
}

// NewIntegerComparable creates sharded rwlock maps with comparable values.
func NewIntegerComparable[K constraints.Integer, V comparable](shardsCount, defaultSize int) GenericComparable[K, V] {
	return NewGenericComparable[K, V](shardsCount, defaultSize, func(key K) int {
		return int(key) % shardsCount
	})
}
