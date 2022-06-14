package smap_test

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"testing"

	"github.com/lispad/go-generics-tools/smap"
)

func BenchmarkIntegerSMap_ConcurrentGet(b *testing.B) {
	sm := smap.NewIntegerComparable[uint16, uint64](smap.HeuristicOptimalDistribution(math.MaxUint16))
	for i := uint16(0); i < math.MaxUint16; i++ {
		sm.Store(i, uint64(i))
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := uint16(rand.Uint32())
		for pb.Next() {
			sm.Load(i)
			i++
		}
	})
}

func BenchmarkSyncMap_ConcurrentGet(b *testing.B) {
	sm := sync.Map{}
	for i := uint16(0); i < math.MaxUint16; i++ {
		sm.Store(i, uint64(i))
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := uint16(rand.Uint32())
		for pb.Next() {
			sm.Load(i)
			i++
		}
	})
}

func BenchmarkLockMap_ConcurrentGet(b *testing.B) {
	sm := make(map[uint16]uint64)
	var mutex sync.RWMutex
	for i := uint16(0); i < math.MaxUint16; i++ {
		sm[i] = uint64(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := uint16(rand.Uint32())
		var (
			val uint64
			ok  bool
		)
		for pb.Next() {
			mutex.RLock()
			val, ok = sm[i]
			mutex.RUnlock()
			i++
		}
		_ = val
		_ = ok
	})
}

func BenchmarkIntegerShardedMap_ConcurrentSet(b *testing.B) {
	sm := smap.NewIntegerComparable[uint16, uint64](smap.HeuristicOptimalDistribution(math.MaxUint16))
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := uint16(rand.Uint32())
		for pb.Next() {
			sm.Store(i, uint64(i))
			i++
		}
	})
}

func BenchmarkSyncMap_ConcurrentSet(b *testing.B) {
	sm := sync.Map{}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := uint16(rand.Uint32())
		for pb.Next() {
			sm.Store(i, uint64(i))
			i++
		}
	})
}

func BenchmarkLockMap_ConcurrentSet(b *testing.B) {
	sm := make(map[uint16]uint64)
	var mutex sync.RWMutex

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := uint16(rand.Uint32())
		for pb.Next() {
			mutex.Lock()
			sm[i] = uint64(i)
			mutex.Unlock()
			i++
		}
	})
}

func BenchmarkIntegerShardedMap_ConcurrentGetSet5(b *testing.B) {
	startMemory, endMemory := prepareMemStats()
	benchmarkIntegerShardedMapConcurrentGetSet(b, 19) // 95% reads, 5% writes
	reportMemStats("sharded rwlocked map", &startMemory, &endMemory)
}

func BenchmarkSyncMap_ConcurrentGetSet5(b *testing.B) {
	startMemory, endMemory := prepareMemStats()
	benchmarkSyncMapConcurrentGetSet(b, 19) // 95% reads, 5% writes
	reportMemStats("sync map", &startMemory, &endMemory)
}

func BenchmarkLockMap_ConcurrentGetSet5(b *testing.B) {
	startMemory, endMemory := prepareMemStats()
	benchmarkLockMapConcurrentGetSet(b, 19) // 95% reads, 5% writes
	reportMemStats("lock map", &startMemory, &endMemory)
}

func BenchmarkIntegerShardedMap_ConcurrentGetSet50(b *testing.B) {
	benchmarkIntegerShardedMapConcurrentGetSet(b, 1) // 50% reads, 50% writes
}

func BenchmarkSyncMap_ConcurrentGetSet50(b *testing.B) {
	benchmarkSyncMapConcurrentGetSet(b, 1) // 50% reads, 50% writes
}

func BenchmarkLockMap_ConcurrentGetSet50(b *testing.B) {
	benchmarkLockMapConcurrentGetSet(b, 1) // 50% reads, 50% writes
}

func BenchmarkIntegerShardedMap_ConcurrentGetSet1(b *testing.B) {
	benchmarkIntegerShardedMapConcurrentGetSet(b, 99) // 99% reads, 1% writes
}

func BenchmarkSyncMap_ConcurrentGetSet1(b *testing.B) {
	benchmarkSyncMapConcurrentGetSet(b, 99) // 99% reads, 1% writes
}

func BenchmarkLockMap_ConcurrentGetSet1(b *testing.B) {
	benchmarkLockMapConcurrentGetSet(b, 99) // 99% reads, 1% writes
}

func benchmarkIntegerShardedMapConcurrentGetSet(b *testing.B, ratio int) {
	sm := smap.NewIntegerComparable[uint16, uint64](smap.HeuristicOptimalDistribution(math.MaxUint16))
	b.ReportAllocs()
	b.ResetTimer()

	benchmarkConcurrentGetSetRatio(b, ratio, func(k uint16) (uint64, bool) {
		return sm.Load(k)
	}, func(k uint16, v uint64) {
		sm.Store(k, v)
	})
}

func benchmarkSyncMapConcurrentGetSet(b *testing.B, ratio int) {
	sm := sync.Map{}
	b.ReportAllocs()
	b.ResetTimer()

	benchmarkConcurrentGetSetRatio(b, ratio, func(k uint16) (uint64, bool) {
		val, ok := sm.Load(k)
		if ok {
			return val.(uint64), ok
		}
		return 0, false
	}, func(k uint16, v uint64) {
		sm.Store(k, v)
	})
}

func benchmarkLockMapConcurrentGetSet(b *testing.B, ratio int) {
	sm := make(map[uint16]uint64, math.MaxUint16)
	var mutex sync.RWMutex
	b.ReportAllocs()
	b.ResetTimer()

	benchmarkConcurrentGetSetRatio(b, ratio, func(k uint16) (uint64, bool) {
		mutex.RLock()
		val, ok := sm[k]
		mutex.RUnlock()
		return val, ok
	}, func(k uint16, v uint64) {
		mutex.Lock()
		sm[k] = v
		mutex.Unlock()
	})
}

func benchmarkConcurrentGetSetRatio(b *testing.B, ratio int, get func(k uint16) (uint64, bool), set func(uint16, uint64)) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := uint16(rand.Uint32())
		j := uint16(rand.Uint32())
		for pb.Next() {
			set(i, uint64(i))
			for k := 0; k < ratio; k++ {
				get(j)
				j++
			}
			i++
		}
	})
}

func prepareMemStats() (startMemory, endMemory runtime.MemStats) {
	runtime.ReadMemStats(&startMemory)
	return
}

func reportMemStats(test string, startMemory, endMemory *runtime.MemStats) {
	runtime.ReadMemStats(endMemory)
	fmt.Printf("Test %s. Memory used %.1fKb\n", test, float32(endMemory.TotalAlloc-startMemory.TotalAlloc)/1024)
}
