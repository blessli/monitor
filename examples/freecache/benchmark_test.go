package main

import (
	"testing"
	"github.com/coocood/freecache"
	"time"
)

// BenchmarkGetAndSet_Freecache-2             73832             16154 ns/op            4384 B/op          9 allocs/op
// BenchmarkGetAndSet_Freecache-2       	  71943             15297 ns/op            4384 B/op          9 allocs/op
// BenchmarkGetAndSet_Freecache-2             75753             15712 ns/op            4384 B/op          9 allocs/op
// BenchmarkGetAndSet_Freecache-2             78002             15478 ns/op            4384 B/op          9 allocs/op
func BenchmarkGetAndSet_Freecache(B *testing.B) {
	B.Logf("B value is %d", B.N)
	for i := 0; i < B.N; i++ {
		key := getRandNum()
		v, err := cacheInstance.Get([]byte(key))
		if err != nil && err.Error() != freecache.ErrNotFound.Error() {
			B.Logf("freecache get error: %v", err)
			return
		}
		if len(v) == 0 {
			if err := cacheInstance.Set([]byte(key), cacheValue, 2); err != nil {
				B.Logf("ser error: %v", err)
			}
		}
	}
}
// BenchmarkGetAndSet_Gocache-2              83656             12639 ns/op             152 B/op          6 allocs/op
// BenchmarkGetAndSet_Gocache-2              94089             12563 ns/op             152 B/op          6 allocs/op
// BenchmarkGetAndSet_Gocache-2              94388             12788 ns/op             152 B/op          6 allocs/op
// BenchmarkGetAndSet_Gocache-2              94209             12924 ns/op             152 B/op          6 allocs/op
func BenchmarkGetAndSet_Gocache(B *testing.B) {
	B.Logf("B value is %d", B.N)
	for i := 0; i < B.N; i++ {
		key := getRandNum()
		_, found := cacheInstance2.Get(key)
		if !found {
			cacheInstance2.Set(key, cacheValue, 2*time.Second)
		}
	}

}
