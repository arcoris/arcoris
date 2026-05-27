// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atomicx

import "testing"

func BenchmarkPaddedUint64Load(b *testing.B) {
	var val PaddedUint64
	val.Store(1)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = val.Load()
	}
}

func BenchmarkPaddedUint64Add(b *testing.B) {
	var val PaddedUint64

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		val.Add(1)
	}
}

func BenchmarkPaddedUint64CompareAndSwap(b *testing.B) {
	var val PaddedUint64

	b.ReportAllocs()
	b.ResetTimer()

	for i := uint64(0); i < uint64(b.N); i++ {
		if !val.CompareAndSwap(i, i+1) {
			b.Fatal("unexpected CAS failure")
		}
	}
}

func BenchmarkUint64CounterInc(b *testing.B) {
	var counter Uint64Counter

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		counter.Inc()
	}
}

func BenchmarkUint64CounterAdd(b *testing.B) {
	var counter Uint64Counter

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		counter.Add(3)
	}
}

func BenchmarkUint64CounterSnapshot(b *testing.B) {
	var counter Uint64Counter
	counter.Add(42)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = counter.Snapshot()
	}
}

func BenchmarkUint64CounterDeltaSince(b *testing.B) {
	prev := Uint64CounterSnapshot{Value: 42}
	cur := Uint64CounterSnapshot{Value: 84}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = cur.DeltaSince(prev)
	}
}

func BenchmarkUint64GaugeAdd(b *testing.B) {
	var gauge Uint64Gauge

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		gauge.Add(1)
	}
}

func BenchmarkUint64GaugeTryAdd(b *testing.B) {
	var gauge Uint64Gauge

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = gauge.TryAdd(1)
	}
}

func BenchmarkUint64GaugeSub(b *testing.B) {
	var gauge Uint64Gauge
	gauge.Store(uint64(b.N))

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		gauge.Sub(1)
	}
}

func BenchmarkUint64GaugeTrySub(b *testing.B) {
	var gauge Uint64Gauge
	gauge.Store(uint64(b.N))

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = gauge.TrySub(1)
	}
}

func BenchmarkInt64GaugeAdd(b *testing.B) {
	var gauge Int64Gauge

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		gauge.Add(1)
	}
}

func BenchmarkInt64GaugeTryAdd(b *testing.B) {
	var gauge Int64Gauge

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = gauge.TryAdd(1)
	}
}

func BenchmarkInt64GaugeSub(b *testing.B) {
	var gauge Int64Gauge
	gauge.Store(int64(b.N))

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		gauge.Sub(1)
	}
}

func BenchmarkInt64GaugeTrySub(b *testing.B) {
	var gauge Int64Gauge
	gauge.Store(int64(b.N))

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = gauge.TrySub(1)
	}
}

func BenchmarkPaddedPointerLoad(b *testing.B) {
	var ptr PaddedPointer[pointerTestValue]
	value := &pointerTestValue{value: 1}
	ptr.Store(value)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = ptr.Load()
	}
}

func BenchmarkPaddedPointerStore(b *testing.B) {
	var ptr PaddedPointer[pointerTestValue]
	value := &pointerTestValue{value: 1}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		ptr.Store(value)
	}
}

func BenchmarkUint64CounterIncParallel(b *testing.B) {
	var counter Uint64Counter

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Inc()
		}
	})
}

func BenchmarkPaddedUint64AddParallel(b *testing.B) {
	var val PaddedUint64

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			val.Add(1)
		}
	})
}

func BenchmarkUint64GaugeAddSubParallel(b *testing.B) {
	var gauge Uint64Gauge

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gauge.Add(1)
			gauge.Sub(1)
		}
	})
}
