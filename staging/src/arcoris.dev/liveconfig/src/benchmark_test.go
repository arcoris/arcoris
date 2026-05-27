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

package liveconfig

import (
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkHolderSnapshot(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = h.Snapshot()
	}
}

func BenchmarkHolderStamped(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = h.Stamped()
	}
}

func BenchmarkHolderRevision(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = h.Revision()
	}
}

func BenchmarkHolderLastError(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = h.LastError()
	}
}

func BenchmarkHolderApplyPublished(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := range b.N {
		_, _ = h.Apply(testConfig{Name: "next", Limit: i})
	}
}

func BenchmarkHolderApplyEqualNoop(b *testing.B) {
	h := newBenchmarkHolder(b, WithEqual(equalTestConfig))

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = h.Apply(testConfig{Name: "initial", Limit: 1})
	}
}

func BenchmarkHolderApplyValidateFailed(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = h.Apply(testConfig{Name: "bad", Limit: -1})
	}
}

func BenchmarkHolderApplyNormalizeFailed(b *testing.B) {
	errNormalize := benchmarkError{}
	h := newBenchmarkHolder(b, WithNormalizer(func(cfg testConfig) (testConfig, error) {
		if cfg.Name == "bad" {
			return testConfig{}, errNormalize
		}
		return cfg, nil
	}))

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = h.Apply(testConfig{Name: "bad", Limit: 1})
	}
}

func BenchmarkHolderSnapshotParallel(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = h.Snapshot()
		}
	})
}

func BenchmarkHolderApplyPublishedParallel(b *testing.B) {
	h := newBenchmarkHolder(b)
	var next atomic.Int64

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limit := int(next.Add(1))
			_, _ = h.Apply(testConfig{Name: "next", Limit: limit})
		}
	})
}

func BenchmarkHolderSnapshotWhileApplying(b *testing.B) {
	h := newBenchmarkHolder(b)
	var next atomic.Int64
	done := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				limit := int(next.Add(1))
				_, _ = h.Apply(testConfig{Name: "next", Limit: limit})
			}
		}
	}()

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = h.Snapshot()
	}

	b.StopTimer()
	close(done)
	wg.Wait()
}

type benchmarkError struct{}

func (benchmarkError) Error() string {
	return "benchmark error"
}

func newBenchmarkHolder(b *testing.B, opts ...Option[testConfig]) *Holder[testConfig] {
	b.Helper()

	base := []Option[testConfig]{
		WithClock[testConfig](newTestClock()),
		WithClone(cloneTestConfig),
		WithValidator(validTestConfig),
	}
	base = append(base, opts...)

	h, err := New(testConfig{Name: "initial", Limit: 1}, base...)
	if err != nil {
		b.Fatalf("New() error = %v", err)
	}
	return h
}
