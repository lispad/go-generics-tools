package smap

import (
	"runtime"
	"sync"
)

// Generic stores data in N shards, with rw mutex for each.
type Generic[K comparable, V any] struct {
	shards        []map[K]V
	locks         []sync.RWMutex
	shardDetector func(key K) int
}

// NewGeneric creates generic RWLocked Sharded map.
// shardDetector should be idempotent function.
func NewGeneric[K comparable, V any](shardsCount, defaultSize int, shardDetector func(key K) int) Generic[K, V] {
	sm := Generic[K, V]{
		shards:        make([]map[K]V, shardsCount),
		locks:         make([]sync.RWMutex, shardsCount),
		shardDetector: shardDetector,
	}
	for i := 0; i < shardsCount; i++ {
		sm.shards[i] = make(map[K]V, defaultSize)
		sm.locks[i] = sync.RWMutex{}
	}
	return sm
}

// Load returns the value stored in the map for a key, or nil if no value is present.
// The ok result indicates whether value was found in the map.
func (sm Generic[K, V]) Load(key K) (V, bool) {
	shardID := sm.shardDetector(key)
	sm.locks[shardID].RLock()
	value, ok := sm.shards[shardID][key]
	sm.locks[shardID].RUnlock()
	return value, ok
}

// Store sets the value for a key.
func (sm Generic[K, V]) Store(key K, value V) {
	shardID := sm.shardDetector(key)
	sm.locks[shardID].Lock()
	sm.shards[shardID][key] = value
	sm.locks[shardID].Unlock()
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (sm Generic[K, V]) LoadAndDelete(key K) (V, bool) {
	shardID := sm.shardDetector(key)
	sm.locks[shardID].Lock()
	value, ok := sm.shards[shardID][key]
	if ok {
		delete(sm.shards[shardID], key)
	}
	sm.locks[shardID].Unlock()
	return value, ok
}

// LoadOrCreate returns the existing value for the key if present.
// Otherwise, it calls generator func, stores and returns the generator's result.
// Generator will not be called if key present.
// The loaded result is true if the value was loaded, false if stored.
func (sm Generic[K, V]) LoadOrCreate(key K, generator func() V) (V, bool) {
	shardID := sm.shardDetector(key)
	sm.locks[shardID].RLock()
	value, ok := sm.shards[shardID][key]
	sm.locks[shardID].RUnlock()
	if ok {
		return value, ok
	}

	sm.locks[shardID].Lock()
	value, ok = sm.shards[shardID][key]
	if !ok {
		value = generator()
		sm.shards[shardID][key] = value
	}
	sm.locks[shardID].Unlock()
	return value, ok
}

// Delete deletes the value for a key.
func (sm Generic[K, V]) Delete(key K) {
	shardID := sm.shardDetector(key)
	sm.locks[shardID].Lock()
	delete(sm.shards[shardID], key)
	sm.locks[shardID].Unlock()
}

// Range calls cb sequentially for each key and value present in the map.
// If cb returns false, range stops the iteration.
//
// Range performs like sync.Range: does not correspond to any consistent snapshot of the Map's contents:
// no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by cb), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even cb itself may call any method on sm.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (sm Generic[K, V]) Range(cb func(K, V) bool) {
	keys := make([]K, 0)
	for i := range sm.locks {
		keys = keys[:0]
		sm.locks[i].RLock()
		for k := range sm.shards[i] {
			keys = append(keys, k)
		}
		sm.locks[i].RUnlock()

		for _, key := range keys {
			sm.locks[i].RLock()
			value, ok := sm.shards[i][key]
			sm.locks[i].RUnlock()
			if ok {
				if !cb(key, value) {
					return
				}
			}
		}
	}
}

// ShardID returns shard number for given key.
func (sm Generic[K, V]) ShardID(key K) int {
	return sm.shardDetector(key)
}

// ShardsCount returns shards count, given on initialisation.
func (sm Generic[K, V]) ShardsCount() int {
	return len(sm.locks)
}

// LockShard locks shard with given id.
// Could be useful with Unblocked* functions. Other calls to sm could be locked.
// Use with caution, only when benchmark shows significant performance changes.
func (sm Generic[K, V]) LockShard(id int) {
	sm.locks[id].Lock()
}

// RLockShard locks for read shard with given id.
// Could be useful with Unblocked* functions. Other calls to sm could be rlocked.
// Use with caution, only when benchmark shows significant performance changes.
func (sm Generic[K, V]) RLockShard(id int) {
	sm.locks[id].RLock()
}

// UnlockShard unlocks shard with given id.
func (sm Generic[K, V]) UnlockShard(id int) {
	sm.locks[id].Unlock()
}

// RUnlockShard unlocks for read shard with given id.
func (sm Generic[K, V]) RUnlockShard(id int) {
	sm.locks[id].RUnlock()
}

// UnblockedGet returns value, without locks.
// Use with caution, only when lock or rlock were taken for shard.
func (sm Generic[K, V]) UnblockedGet(key K) (V, bool) {
	value, ok := sm.shards[sm.shardDetector(key)][key]
	return value, ok
}

// UnblockedSet sets value, without locks.
// Use with caution, only when lock were taken for shard.
func (sm Generic[K, V]) UnblockedSet(key K, value V) {
	sm.shards[sm.shardDetector(key)][key] = value
}

// UnblockedShardRange calls cb sequentially for each key and value present in the maps shard.
// Use with caution, only when lock or rlock were taken for shard.
func (sm Generic[K, V]) UnblockedShardRange(shardID int, cb func(key K, value V) bool) {
	for key, value := range sm.shards[shardID] {
		if !cb(key, value) {
			break
		}
	}
}

// HeuristicOptimalShardsCount returns shards count.
// Use with caution, It's point for change, and could be changed in the future.
func HeuristicOptimalShardsCount() int {
	procs := runtime.GOMAXPROCS(-1)
	return procs * procs * 24
}

// HeuristicOptimalDistribution returns shards count, and shard size for given map size.
// Use with caution, It's point for change, and could be changed in the future.
func HeuristicOptimalDistribution(expectedItemsCount int) (int, int) {
	shards := HeuristicOptimalShardsCount()
	shardSize := expectedItemsCount/shards + 1
	return shards, shardSize
}
