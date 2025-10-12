package smap

// GenericComparable stores data in N shards, with rw mutex for each.
// Additional CompareAndSwap method added for comparable values.
type GenericComparable[K comparable, V comparable] struct {
	Generic[K, V]
}

// NewGenericComparable creates generic RWLocked Sharded map for comparable values.
// shardDetector should be idempotent function.
func NewGenericComparable[K comparable, V comparable](
	shardsCount int,
	defaultSize int,
	shardDetector func(key K) int,
) GenericComparable[K, V] {
	return GenericComparable[K, V]{
		Generic: NewGeneric[K, V](shardsCount, defaultSize, shardDetector),
	}
}

// CompareAndSwap executes the compare-and-swap operation for the Key & Value pair.
// If and only if key exists, and value for key equals old, value will be changed to actual.
// Otherwise, returns current value.
// The ok result indicates whether value was changed to actual in the map.
func (sm GenericComparable[K, V]) CompareAndSwap(key K, old, actual V) (V, bool) {
	shardID := sm.shardDetector(key)
	sm.locks[shardID].Lock()
	current, ok := sm.shards[shardID][key]
	if ok && current == old {
		sm.shards[shardID][key] = actual
		sm.locks[shardID].Unlock()
		return actual, true
	}

	sm.locks[shardID].Unlock()
	return current, false
}
