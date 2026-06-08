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
	"errors"
	"sync"
	"testing"

	"arcoris.dev/snapshot"
)

func TestSnapshotReturnsCurrentValue(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1})

	_, err := h.Apply(testConfig{Name: "current", Limit: 2})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	snap := h.Snapshot()
	if got, want := snap.Value.Name, "current"; got != want {
		t.Fatalf("Snapshot().Value.Name = %q, want %q", got, want)
	}
	if got, want := snap.Value.Limit, 2; got != want {
		t.Fatalf("Snapshot().Value.Limit = %d, want %d", got, want)
	}
}

func TestStampedReturnsPublishedTimestamp(t *testing.T) {
	clk := newTestClock()
	h, err := New(
		testConfig{Name: "initial"},
		WithClock[testConfig](clk),
		WithValidator(validTestConfig),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	stamped := h.Stamped()
	if stamped.IsZeroRevision() {
		t.Fatal("Stamped().Revision is zero")
	}
	if stamped.Updated.IsZero() {
		t.Fatal("Stamped().Updated is zero")
	}
	if !stamped.Updated.Equal(clk.now) {
		t.Fatalf("Stamped().Updated = %s, want %s", stamped.Updated, clk.now)
	}
	if got, want := stamped.Value.Name, "initial"; got != want {
		t.Fatalf("Stamped value name = %q, want %q", got, want)
	}
}

func TestRevisionReturnsCurrentRevision(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial"})

	if got, want := h.Revision(), h.Snapshot().Revision; got != want {
		t.Fatalf("Revision() = %d, want %d", got, want)
	}
}

func TestLastErrorRecordsNormalizeFailure(t *testing.T) {
	errNormalize := errors.New("normalize failed")
	h := newTestHolder(t, testConfig{Name: "initial"}, WithNormalizer(func(cfg testConfig) (testConfig, error) {
		if cfg.Name == "bad" {
			return testConfig{}, errNormalize
		}
		return cfg, nil
	}))

	_, err := h.Apply(testConfig{Name: "bad"})
	if err != errNormalize {
		t.Fatalf("Apply() error = %v, want %v", err, errNormalize)
	}
	if got := h.LastError(); got != errNormalize {
		t.Fatalf("LastError() = %v, want %v", got, errNormalize)
	}
}

func TestLastErrorRecordsValidateFailure(t *testing.T) {
	errValidate := errors.New("validate failed")
	h := newTestHolder(t, testConfig{Name: "initial"}, WithValidator(func(cfg testConfig) error {
		if cfg.Name == "bad" {
			return errValidate
		}
		return validTestConfig(cfg)
	}))

	_, err := h.Apply(testConfig{Name: "bad"})
	if err != errValidate {
		t.Fatalf("Apply() error = %v, want %v", err, errValidate)
	}
	if got := h.LastError(); got != errValidate {
		t.Fatalf("LastError() = %v, want %v", got, errValidate)
	}
}

func TestLastErrorIsOverwrittenByNewerRejectedApply(t *testing.T) {
	errNormalize := errors.New("normalize failed")
	errValidate := errors.New("validate failed")
	h := newTestHolder(
		t,
		testConfig{Name: "initial"},
		WithNormalizer(func(cfg testConfig) (testConfig, error) {
			if cfg.Name == "normalize" {
				return testConfig{}, errNormalize
			}
			return cfg, nil
		}),
		WithValidator(func(cfg testConfig) error {
			if cfg.Name == "validate" {
				return errValidate
			}
			return nil
		}),
	)

	_, _ = h.Apply(testConfig{Name: "normalize"})
	_, _ = h.Apply(testConfig{Name: "validate"})
	if got := h.LastError(); got != errValidate {
		t.Fatalf("LastError() = %v, want %v", got, errValidate)
	}
}

func TestLastErrorClearedByPublishedApply(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1})
	_, _ = h.Apply(testConfig{Name: "bad", Limit: -1})

	_, err := h.Apply(testConfig{Name: "good", Limit: 2})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if got := h.LastError(); got != nil {
		t.Fatalf("LastError() = %v, want nil", got)
	}
}

func TestLastErrorClearedByEqualApply(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1}, WithEqual(equalTestConfig))
	_, _ = h.Apply(testConfig{Name: "bad", Limit: -1})

	_, err := h.Apply(testConfig{Name: "initial", Limit: 1})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if got := h.LastError(); got != nil {
		t.Fatalf("LastError() = %v, want nil", got)
	}
}

func TestReadMethodsDoNotClearLastError(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1})
	_, err := h.Apply(testConfig{Name: "bad", Limit: -1})
	if err == nil {
		t.Fatal("Apply() error = nil, want validation error")
	}

	_ = h.Snapshot()
	_ = h.Stamped()
	_ = h.Revision()
	_ = h.LastError()
	if got := h.LastError(); got != err {
		t.Fatalf("LastError() = %v, want %v", got, err)
	}
}

func TestHolderSnapshotDoesNotClonePublishedValue(t *testing.T) {
	h, err := New(mutableConfig{Tags: []string{"a"}}, WithClone(cloneMutableConfig))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	snap := h.Snapshot()
	snap.Value.Tags[0] = "changed"
	if got, want := h.Snapshot().Value.Tags[0], "changed"; got != want {
		t.Fatalf("Snapshot().Value.Tags[0] = %q, want %q", got, want)
	}
}

func TestHolderStampedDoesNotClonePublishedValue(t *testing.T) {
	h, err := New(mutableConfig{Nested: &nestedConfig{Values: []string{"a"}}}, WithClone(cloneMutableConfig))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	stamped := h.Stamped()
	stamped.Value.Nested.Values[0] = "changed"
	if got, want := h.Stamped().Value.Nested.Values[0], "changed"; got != want {
		t.Fatalf("Stamped().Value.Nested.Values[0] = %q, want %q", got, want)
	}
}

func TestConcurrentSnapshotAndApply(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 0})

	const writers = 32
	const readers = 32
	const readsPerReader = 64

	start := make(chan struct{})
	errs := make(chan error, writers)
	var wg sync.WaitGroup
	wg.Add(writers + readers)

	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < readsPerReader; j++ {
				_ = h.Snapshot()
				_ = h.Stamped()
				_ = h.Revision()
				_ = h.LastError()
			}
		}()
	}

	for i := 0; i < writers; i++ {
		i := i
		go func() {
			defer wg.Done()
			<-start
			_, err := h.Apply(testConfig{Name: "writer", Limit: i})
			if err != nil {
				errs <- err
			}
		}()
	}

	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		t.Fatalf("Apply() error = %v", err)
	}

	want := snapshot.Revision(1 + writers)
	if got := h.Revision(); got != want {
		t.Fatalf("Revision() = %d, want %d", got, want)
	}
}

func TestConcurrentApplyAndLastErrorIsRaceFree(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1})

	runConcurrent(64, func(i int) {
		if i%2 == 0 {
			_, _ = h.Apply(testConfig{Name: "bad", Limit: -1})
			return
		}
		_ = h.LastError()
		_, _ = h.Apply(testConfig{Name: "good", Limit: i})
		_ = h.LastError()
	})

	if h.Snapshot().Value.Limit < 0 {
		t.Fatal("rejected candidate became visible")
	}
}

func BenchmarkHolderSnapshot(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		benchmarkSnapshotSink = h.Snapshot()
	}
}

func BenchmarkHolderStamped(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		benchmarkStampedSink = h.Stamped()
	}
}

func BenchmarkHolderRevision(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		benchmarkRevisionSink = h.Revision()
	}
}

func BenchmarkHolderLastError(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		benchmarkErrorSink = h.LastError()
	}
}

func BenchmarkHolderSnapshotParallel(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var local snapshotLocal
		for pb.Next() {
			local.snapshot = h.Snapshot()
		}
		local.keep()
	})
}

func BenchmarkHolderStampedParallel(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var local stampedLocal
		for pb.Next() {
			local.stamped = h.Stamped()
		}
		local.keep()
	})
}

func BenchmarkHolderRevisionParallel(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var local revisionLocal
		for pb.Next() {
			local.revision = h.Revision()
		}
		local.keep()
	})
}

func BenchmarkHolderLastErrorParallel(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var local errorLocal
		for pb.Next() {
			local.err = h.LastError()
		}
		local.keep()
	})
}

func BenchmarkHolderSnapshotWhileApplying(b *testing.B) {
	h := newBenchmarkHolder(b)
	stop := startBenchmarkApplier(h)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		benchmarkSnapshotSink = h.Snapshot()
	}

	b.StopTimer()
	stop()
}

func BenchmarkHolderRevisionWhileApplying(b *testing.B) {
	h := newBenchmarkHolder(b)
	stop := startBenchmarkApplier(h)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		benchmarkRevisionSink = h.Revision()
	}

	b.StopTimer()
	stop()
}

func BenchmarkHolderLastErrorWhileApplying(b *testing.B) {
	h := newBenchmarkHolder(b)
	stop := startBenchmarkApplier(h)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		benchmarkErrorSink = h.LastError()
	}

	b.StopTimer()
	stop()
}
