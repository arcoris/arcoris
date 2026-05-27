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

package snapshot

import (
	"sync/atomic"
	"testing"
)

func BenchmarkPublisherSnapshotSmallValue(b *testing.B) {
	publisher := NewPublisher[int]()
	publisher.Publish(42)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = publisher.Snapshot()
	}
}

func BenchmarkPublisherSnapshotParallel(b *testing.B) {
	publisher := NewPublisher[int]()
	publisher.Publish(42)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = publisher.Snapshot()
		}
	})
}

func BenchmarkPublisherPublishSmallValue(b *testing.B) {
	publisher := NewPublisher[int]()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = publisher.Publish(i)
	}
}

func BenchmarkPublisherPublishParallel(b *testing.B) {
	publisher := NewPublisher[int]()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = publisher.Publish(1)
		}
	})
}

func BenchmarkPublisherSnapshotWhilePublishing(b *testing.B) {
	publisher := NewPublisher[int]()
	publisher.Publish(0)

	var stop atomic.Bool
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; !stop.Load(); i++ {
			_ = publisher.Publish(i)
		}
	}()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = publisher.Snapshot()
		}
	})

	b.StopTimer()
	stop.Store(true)
	<-done
}

func BenchmarkPublisherSnapshotSlice100(b *testing.B) {
	publisher := NewPublisher[[]string]()
	publisher.Publish(make([]string, 100))

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = publisher.Snapshot()
	}
}

func BenchmarkPublisherPublishSlice100(b *testing.B) {
	publisher := NewPublisher[[]string]()
	val := make([]string, 100)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = publisher.Publish(val)
	}
}
