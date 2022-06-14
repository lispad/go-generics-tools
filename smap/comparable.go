package smap

// GenericComparable stores data in N shards, with rw mutex for each.
// Additional CompareAndSwap method added for comparable values.
type GenericComparable[K comparable, V comparable] struct {
	Generic[K, V]
}

// NewGenericComparable creates generic RWLocked Sharded map for comparable values.
// shardDetector should be idempotent function.
func NewGenericComparable[K comparable, V comparable](shardsCount, defaultSize int, shardDetector func(key K) int) GenericComparable[K, V] {
	return GenericComparable[K, V]{
		Generic: NewGeneric[K, V](shardsCount, defaultSize, shardDetector),
	}
}

// CompareAndSwap executes the compare-and-swap operation for the Key & Value pair.
// If and only if key exists, and value for key equals old, value will be changed to new.
// Otherwise, returns current value.
// The ok result indicates whether value was changed to new in the map.
func (sm GenericComparable[K, V]) CompareAndSwap(key K, old, new V) (V, bool) {
	shardID := sm.shardDetector(key)
	sm.locks[shardID].Lock()
	if current, ok := sm.shards[shardID][key]; ok && current == old {
		sm.shards[shardID][key] = new
		sm.locks[shardID].Unlock()
		return new, true
	} else {
		sm.locks[shardID].Unlock()
		return current, false
	}
}
