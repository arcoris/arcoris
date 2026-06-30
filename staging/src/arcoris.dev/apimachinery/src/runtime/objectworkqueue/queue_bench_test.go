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

package objectworkqueue

import (
	"context"
	"sync/atomic"
	"testing"
)

// Benchmarks are baseline measurements for the current one-mutex, map + list
// implementation. Compare results with benchstat before changing data
// structures.

func BenchmarkTryAddDistinct(b *testing.B) {
	queue := newBenchQueue(b, b.N)
	items := benchItems(b.N)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := queue.TryAdd(items[i]); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTryAddDuplicateQueued(b *testing.B) {
	queue := newBenchQueue(b, 1)
	item := testItem(1)
	if err := queue.TryAdd(item); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := queue.TryAdd(item); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTryAddDuplicateProcessing(b *testing.B) {
	queue := newBenchQueue(b, 1)
	item := testItem(1)
	if err := queue.TryAdd(item); err != nil {
		b.Fatal(err)
	}
	if _, err := queue.Get(context.Background()); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := queue.TryAdd(item); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAddGetDoneSerial(b *testing.B) {
	queue := newBenchQueue(b, 1)
	item := testItem(1)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := queue.Add(ctx, item); err != nil {
			b.Fatal(err)
		}
		got, err := queue.Get(ctx)
		if err != nil {
			b.Fatal(err)
		}
		if err := queue.Done(got); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDirtyRequeue(b *testing.B) {
	queue := newBenchQueue(b, 1)
	item := testItem(1)
	ctx := context.Background()
	if err := queue.Add(ctx, item); err != nil {
		b.Fatal(err)
	}
	if _, err := queue.Get(ctx); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := queue.Add(ctx, item); err != nil {
			b.Fatal(err)
		}
		if err := queue.Done(item); err != nil {
			b.Fatal(err)
		}
		if _, err := queue.Get(ctx); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParallelTryAddDuplicate(b *testing.B) {
	queue := newBenchQueue(b, 1)
	item := testItem(1)
	if err := queue.TryAdd(item); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mustBench(queue.TryAdd(item))
		}
	})
}

func BenchmarkParallelTryAddDistinct(b *testing.B) {
	queue := newBenchQueue(b, b.N)
	items := benchItems(b.N)
	var next atomic.Uint64

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := int(next.Add(1)) - 1
			mustBench(queue.TryAdd(items[i]))
		}
	})
}

func BenchmarkParallelAddGetDone(b *testing.B) {
	queue := newBenchQueue(b, b.N)
	items := benchItems(b.N)
	ctx := context.Background()
	var next atomic.Uint64

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := int(next.Add(1)) - 1
			mustBench(queue.Add(ctx, items[i]))
			got, err := queue.Get(ctx)
			mustBench(err)
			mustBench(queue.Done(got))
		}
	})
}

func BenchmarkParallelGetDone(b *testing.B) {
	queue := newBenchQueue(b, b.N)
	items := benchItems(b.N)
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		if err := queue.TryAdd(items[i]); err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			got, err := queue.Get(ctx)
			mustBench(err)
			mustBench(queue.Done(got))
		}
	})
}

func BenchmarkStats(b *testing.B) {
	queue := newBenchQueue(b, 1024)
	for i := 0; i < 512; i++ {
		if err := queue.TryAdd(testItem(i)); err != nil {
			b.Fatal(err)
		}
	}
	for i := 0; i < 256; i++ {
		if _, err := queue.Get(context.Background()); err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = queue.Stats()
	}
}

func BenchmarkValidateItem(b *testing.B) {
	item := testItem(1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := validateItem(item); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkKeyForItem(b *testing.B) {
	item := testItem(1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = keyForItem(item)
	}
}

func newBenchQueue(b *testing.B, capacity int) *Queue {
	b.Helper()

	if capacity <= 0 {
		capacity = 1
	}
	queue, err := New(Options{Capacity: capacity})
	if err != nil {
		b.Fatal(err)
	}
	return queue
}

func benchItems(count int) []Item {
	items := make([]Item, count)
	for i := range items {
		items[i] = testItem(i)
	}
	return items
}

func mustBench(err error) {
	if err != nil {
		panic(err)
	}
}
