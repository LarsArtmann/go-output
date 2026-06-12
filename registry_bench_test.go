package output

import (
	"sync"
	"testing"
)

func BenchmarkFormatRegistry_Get(b *testing.B) {
	reg := newFormatRegistry[int]()
	for _, f := range AllFormats {
		reg.register(f, 42)
	}

	b.ResetTimer()

	for b.Loop() {
		reg.get(FormatJSON)
	}
}

func BenchmarkRawMap_Get(b *testing.B) {
	m := make(map[Format]int)
	for _, f := range AllFormats {
		m[f] = 42
	}

	var mu sync.RWMutex

	b.ResetTimer()

	for b.Loop() {
		mu.RLock()

		_ = m[FormatJSON]

		mu.RUnlock()
	}
}

func BenchmarkFormatRegistry_Get_Miss(b *testing.B) {
	reg := newFormatRegistry[int]()

	b.ResetTimer()

	for b.Loop() {
		reg.get(Format("nonexistent"))
	}
}

func BenchmarkRawMap_Get_Miss(b *testing.B) {
	m := make(map[Format]int)

	var mu sync.RWMutex

	b.ResetTimer()

	for b.Loop() {
		mu.RLock()

		_ = m[Format("nonexistent")]

		mu.RUnlock()
	}
}
