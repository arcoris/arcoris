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

func BenchmarkPublisherZeroSnapshot(b *testing.B) {
	var publisher Publisher[int]

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkIntSnapshotSink = publisher.Snapshot()
	}
}

func BenchmarkPublisherSnapshotSmallValue(b *testing.B) {
	publisher := NewPublisher[int]()
	publisher.Publish(42)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkIntSnapshotSink = publisher.Snapshot()
	}
}

func BenchmarkPublisherStampedSmallValue(b *testing.B) {
	publisher := NewPublisher[int]()
	publisher.Publish(42)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkIntStampedSink = publisher.Stamped()
	}
}

func BenchmarkPublisherRevision(b *testing.B) {
	publisher := NewPublisher[int]()
	publisher.Publish(42)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkRevisionSink = publisher.Revision()
	}
}

func BenchmarkPublisherSnapshotParallel(b *testing.B) {
	publisher := NewPublisher[int]()
	publisher.Publish(42)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var snap Snapshot[int]
		for pb.Next() {
			snap = publisher.Snapshot()
		}
		benchmarkSinkMu.Lock()
		benchmarkIntSnapshotSink = snap
		benchmarkSinkMu.Unlock()
	})
}

func BenchmarkPublisherStampedParallel(b *testing.B) {
	publisher := NewPublisher[int]()
	publisher.Publish(42)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var stamped Stamped[int]
		for pb.Next() {
			stamped = publisher.Stamped()
		}
		benchmarkSinkMu.Lock()
		benchmarkIntStampedSink = stamped
		benchmarkSinkMu.Unlock()
	})
}

func BenchmarkPublisherPublishSmallValue(b *testing.B) {
	publisher := NewPublisher[int]()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkIntSnapshotSink = publisher.Publish(i)
	}
}

func BenchmarkPublisherPublishStampedSmallValue(b *testing.B) {
	publisher := NewPublisher[int]()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkIntStampedSink = publisher.PublishStamped(i)
	}
}

func BenchmarkPublisherPublishParallel(b *testing.B) {
	publisher := NewPublisher[int]()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var snap Snapshot[int]
		for pb.Next() {
			snap = publisher.Publish(1)
		}
		benchmarkSinkMu.Lock()
		benchmarkIntSnapshotSink = snap
		benchmarkSinkMu.Unlock()
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
		var snap Snapshot[int]
		for pb.Next() {
			snap = publisher.Snapshot()
		}
		benchmarkSinkMu.Lock()
		benchmarkIntSnapshotSink = snap
		benchmarkSinkMu.Unlock()
	})

	b.StopTimer()
	stop.Store(true)
	<-done
}

func BenchmarkPublisherRevisionWhilePublishing(b *testing.B) {
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
		var rev Revision
		for pb.Next() {
			rev = publisher.Revision()
		}
		benchmarkSinkMu.Lock()
		benchmarkRevisionSink = rev
		benchmarkSinkMu.Unlock()
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
		benchmarkSliceSnapshotSink = publisher.Snapshot()
	}
}

func BenchmarkPublisherPublishSlice100(b *testing.B) {
	publisher := NewPublisher[[]string]()
	val := make([]string, 100)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkSliceSnapshotSink = publisher.Publish(val)
	}
}
