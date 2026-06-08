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
	"sync/atomic"
	"testing"

	"arcoris.dev/snapshot"
)

func TestApplyChangeInvariantMatrix(t *testing.T) {
	errNormalize := errors.New("normalize failed")
	errValidate := errors.New("validate failed")

	tests := []struct {
		name      string
		holder    func(t *testing.T) *Holder[testConfig]
		candidate testConfig
		wantErr   error
		want      ChangeReason
		changed   bool
		current   testConfig
	}{
		{
			name:      "published",
			holder:    func(t *testing.T) *Holder[testConfig] { return newTestHolder(t, testConfig{Name: "initial", Limit: 1}) },
			candidate: testConfig{Name: "next", Limit: 2},
			want:      ChangeReasonPublished,
			changed:   true,
			current:   testConfig{Name: "next", Limit: 2},
		},
		{
			name: "equal",
			holder: func(t *testing.T) *Holder[testConfig] {
				return newTestHolder(t, testConfig{Name: "initial", Limit: 1}, WithEqual(equalTestConfig))
			},
			candidate: testConfig{Name: "initial", Limit: 1},
			want:      ChangeReasonEqual,
			current:   testConfig{Name: "initial", Limit: 1},
		},
		{
			name: "normalize failed",
			holder: func(t *testing.T) *Holder[testConfig] {
				return newTestHolder(t, testConfig{Name: "initial", Limit: 1}, WithNormalizer(func(cfg testConfig) (testConfig, error) {
					if cfg.Name == "bad" {
						return testConfig{}, errNormalize
					}
					return cfg, nil
				}))
			},
			candidate: testConfig{Name: "bad", Limit: 2},
			wantErr:   errNormalize,
			want:      ChangeReasonNormalizeFailed,
			current:   testConfig{Name: "initial", Limit: 1},
		},
		{
			name: "validate failed",
			holder: func(t *testing.T) *Holder[testConfig] {
				return newTestHolder(t, testConfig{Name: "initial", Limit: 1}, WithValidator(func(cfg testConfig) error {
					if cfg.Name == "bad" {
						return errValidate
					}
					return validTestConfig(cfg)
				}))
			},
			candidate: testConfig{Name: "bad", Limit: 2},
			wantErr:   errValidate,
			want:      ChangeReasonValidateFailed,
			current:   testConfig{Name: "initial", Limit: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.holder(t)
			prev := h.Snapshot()

			change, err := h.Apply(tt.candidate)
			if err != tt.wantErr {
				t.Fatalf("Apply() error = %v, want %v", err, tt.wantErr)
			}
			if change.Changed != tt.changed {
				t.Fatalf("Change.Changed = %v, want %v", change.Changed, tt.changed)
			}
			if change.IsChanged() != tt.changed {
				t.Fatalf("Change.IsChanged() = %v, want %v", change.IsChanged(), tt.changed)
			}
			if change.IsNoop() == tt.changed {
				t.Fatalf("Change.IsNoop() = %v, want %v", change.IsNoop(), !tt.changed)
			}
			if got, want := change.Accepted(), tt.want.Accepted(); got != want {
				t.Fatalf("Change.Accepted() = %v, want %v", got, want)
			}
			if got, want := change.Rejected(), tt.want.Rejected(); got != want {
				t.Fatalf("Change.Rejected() = %v, want %v", got, want)
			}
			if got := change.Reason; got != tt.want {
				t.Fatalf("Change.Reason = %s, want %s", got, tt.want)
			}
			if change.Previous.Revision != prev.Revision {
				t.Fatalf("Previous.Revision = %d, want %d", change.Previous.Revision, prev.Revision)
			}
			if tt.changed {
				if got, want := change.Current.Revision, change.Previous.Revision.Next(); got != want {
					t.Fatalf("Current.Revision = %d, want %d", got, want)
				}
			} else if change.Current.Revision != change.Previous.Revision {
				t.Fatalf("Current.Revision = %d, want previous revision %d", change.Current.Revision, change.Previous.Revision)
			}
			if got := h.Revision(); got != change.Current.Revision {
				t.Fatalf("Revision() = %d, want %d", got, change.Current.Revision)
			}
			assertTestSnapshot(t, h.Snapshot(), change.Current)
			if !equalTestConfig(change.Current.Value, tt.current) {
				t.Fatalf("Current.Value = %#v, want %#v", change.Current.Value, tt.current)
			}
			if got := h.LastError(); got != tt.wantErr {
				t.Fatalf("LastError() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestApplyPublishesValidConfig(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1})
	prev := h.Snapshot()

	change, err := h.Apply(testConfig{Name: "next", Limit: 2})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if !change.Changed {
		t.Fatal("Apply().Changed = false, want true")
	}
	if got, want := change.Reason, ChangeReasonPublished; got != want {
		t.Fatalf("Apply().Reason = %s, want %s", got, want)
	}
	if !change.Accepted() {
		t.Fatal("Apply().Accepted() = false, want true")
	}
	if change.Previous.Revision != prev.Revision {
		t.Fatalf("Previous revision = %d, want %d", change.Previous.Revision, prev.Revision)
	}
	if got, want := change.Previous.Value.Name, "initial"; got != want {
		t.Fatalf("Previous value name = %q, want %q", got, want)
	}
	if got, want := change.Current.Revision, prev.Revision.Next(); got != want {
		t.Fatalf("Current revision = %d, want %d", got, want)
	}
	if got, want := h.Revision(), change.Current.Revision; got != want {
		t.Fatalf("Revision() = %d, want %d", got, want)
	}
	if got, want := h.Snapshot().Value.Name, "next"; got != want {
		t.Fatalf("current name = %q, want %q", got, want)
	}
	if got := h.LastError(); got != nil {
		t.Fatalf("LastError() = %v, want nil", got)
	}
}

func TestApplyRejectsInvalidConfigAndKeepsLastGood(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1})
	prev := h.Snapshot()

	change, err := h.Apply(testConfig{Name: "bad", Limit: -1})
	if err == nil {
		t.Fatal("Apply() error = nil, want validation error")
	}
	if change.Changed {
		t.Fatal("Apply().Changed = true, want false")
	}
	if got, want := change.Reason, ChangeReasonValidateFailed; got != want {
		t.Fatalf("Apply().Reason = %s, want %s", got, want)
	}
	if !change.Rejected() {
		t.Fatal("Apply().Rejected() = false, want true")
	}
	if change.Previous.Revision != prev.Revision {
		t.Fatalf("Previous revision = %d, want %d", change.Previous.Revision, prev.Revision)
	}
	if change.Current.Revision != prev.Revision {
		t.Fatalf("Current revision = %d, want %d", change.Current.Revision, prev.Revision)
	}
	if got, want := h.Snapshot().Value.Name, "initial"; got != want {
		t.Fatalf("current name = %q, want %q", got, want)
	}
	if got := h.LastError(); got != err {
		t.Fatalf("LastError() = %v, want %v", got, err)
	}
}

func TestApplyClearsLastErrorAfterSuccessfulApply(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1})

	_, err := h.Apply(testConfig{Name: "bad", Limit: -1})
	if err == nil {
		t.Fatal("Apply() error = nil, want validation error")
	}
	if h.LastError() == nil {
		t.Fatal("LastError() = nil after rejected Apply")
	}

	change, err := h.Apply(testConfig{Name: "good", Limit: 2})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if got, want := change.Reason, ChangeReasonPublished; got != want {
		t.Fatalf("Apply().Reason = %s, want %s", got, want)
	}
	if got := h.LastError(); got != nil {
		t.Fatalf("LastError() = %v, want nil", got)
	}
}

func TestApplyClearsLastErrorAfterSuccessfulNoop(t *testing.T) {
	h := newTestHolder(
		t,
		testConfig{Name: "initial", Limit: 1},
		WithEqual(equalTestConfig),
	)

	_, err := h.Apply(testConfig{Name: "bad", Limit: -1})
	if err == nil {
		t.Fatal("Apply() error = nil, want validation error")
	}
	if h.LastError() == nil {
		t.Fatal("LastError() = nil after rejected Apply")
	}

	change, err := h.Apply(testConfig{Name: "initial", Limit: 1})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if change.Changed {
		t.Fatal("Apply().Changed = true, want false")
	}
	if got, want := change.Reason, ChangeReasonEqual; got != want {
		t.Fatalf("Apply().Reason = %s, want %s", got, want)
	}
	if !change.Accepted() {
		t.Fatal("Apply().Accepted() = false, want true")
	}
	if got := h.LastError(); got != nil {
		t.Fatalf("LastError() = %v, want nil", got)
	}
}

func TestApplyEqualConfigDoesNotPublish(t *testing.T) {
	h := newTestHolder(
		t,
		testConfig{Name: "initial", Limit: 1, Tags: []string{"a"}},
		WithEqual(equalTestConfig),
	)
	prev := h.Snapshot()

	change, err := h.Apply(testConfig{Name: "initial", Limit: 1, Tags: []string{"a"}})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if change.Changed {
		t.Fatal("Apply().Changed = true, want false")
	}
	if got, want := change.Reason, ChangeReasonEqual; got != want {
		t.Fatalf("Apply().Reason = %s, want %s", got, want)
	}
	if change.Previous.Revision != prev.Revision {
		t.Fatalf("Previous revision = %d, want %d", change.Previous.Revision, prev.Revision)
	}
	if change.Current.Revision != prev.Revision {
		t.Fatalf("Current revision = %d, want %d", change.Current.Revision, prev.Revision)
	}
	if got := h.LastError(); got != nil {
		t.Fatalf("LastError() = %v, want nil", got)
	}
}

func TestApplyWithoutEqualPublishesEquivalentConfig(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1})
	prev := h.Snapshot()

	change, err := h.Apply(testConfig{Name: "initial", Limit: 1})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if !change.Changed {
		t.Fatal("Apply().Changed = false, want true")
	}
	if got, want := change.Reason, ChangeReasonPublished; got != want {
		t.Fatalf("Apply().Reason = %s, want %s", got, want)
	}
	if got, want := change.Current.Revision, prev.Revision.Next(); got != want {
		t.Fatalf("Current revision = %d, want %d", got, want)
	}
}

func TestApplyNormalizerErrorReturnsNormalizeFailedReason(t *testing.T) {
	errNormalize := errors.New("normalize failed")
	h := newTestHolder(
		t,
		testConfig{Name: "initial", Limit: 1},
		WithNormalizer(func(cfg testConfig) (testConfig, error) {
			if cfg.Name == "bad" {
				return testConfig{}, errNormalize
			}
			return cfg, nil
		}),
	)
	prev := h.Snapshot()

	change, err := h.Apply(testConfig{Name: "bad", Limit: 2})
	if err != errNormalize {
		t.Fatalf("Apply() error = %v, want %v", err, errNormalize)
	}
	if got, want := change.Reason, ChangeReasonNormalizeFailed; got != want {
		t.Fatalf("Apply().Reason = %s, want %s", got, want)
	}
	if change.Changed {
		t.Fatal("Apply().Changed = true, want false")
	}
	if change.Current.Revision != prev.Revision {
		t.Fatalf("Current revision = %d, want %d", change.Current.Revision, prev.Revision)
	}
	if got := h.LastError(); got != errNormalize {
		t.Fatalf("LastError() = %v, want %v", got, errNormalize)
	}
}

func TestApplyPipelinePanicPreservesLastGoodAndLastError(t *testing.T) {
	errPrevious := errors.New("previous validation error")
	panicCandidate := testConfig{Name: "panic", Limit: 1}

	tests := []struct {
		name   string
		holder func(t *testing.T) *Holder[testConfig]
	}{
		{
			name: "clone",
			holder: func(t *testing.T) *Holder[testConfig] {
				return newTestHolder(
					t,
					testConfig{Name: "initial", Limit: 1},
					WithClone(func(cfg testConfig) testConfig {
						if cfg.Name == "panic" {
							panic("clone panic")
						}
						return cloneTestConfig(cfg)
					}),
					WithValidator(func(cfg testConfig) error {
						if cfg.Name == "previous" {
							return errPrevious
						}
						return validTestConfig(cfg)
					}),
				)
			},
		},
		{
			name: "normalizer",
			holder: func(t *testing.T) *Holder[testConfig] {
				return newTestHolder(
					t,
					testConfig{Name: "initial", Limit: 1},
					WithNormalizer(func(cfg testConfig) (testConfig, error) {
						if cfg.Name == "panic" {
							panic("normalizer panic")
						}
						return cfg, nil
					}),
					WithValidator(func(cfg testConfig) error {
						if cfg.Name == "previous" {
							return errPrevious
						}
						return validTestConfig(cfg)
					}),
				)
			},
		},
		{
			name: "validator",
			holder: func(t *testing.T) *Holder[testConfig] {
				return newTestHolder(
					t,
					testConfig{Name: "initial", Limit: 1},
					WithValidator(func(cfg testConfig) error {
						switch cfg.Name {
						case "panic":
							panic("validator panic")
						case "previous":
							return errPrevious
						default:
							return validTestConfig(cfg)
						}
					}),
				)
			},
		},
		{
			name: "equal",
			holder: func(t *testing.T) *Holder[testConfig] {
				return newTestHolder(
					t,
					testConfig{Name: "initial", Limit: 1},
					WithValidator(func(cfg testConfig) error {
						if cfg.Name == "previous" {
							return errPrevious
						}
						return validTestConfig(cfg)
					}),
					WithEqual(func(_, candidate testConfig) bool {
						if candidate.Name == "panic" {
							panic("equal panic")
						}
						return false
					}),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.holder(t)
			_, err := h.Apply(testConfig{Name: "previous", Limit: 1})
			if err != errPrevious {
				t.Fatalf("Apply() previous error = %v, want %v", err, errPrevious)
			}

			prev := h.Snapshot()
			prevErr := h.LastError()

			requirePanic(t, func() {
				_, _ = h.Apply(panicCandidate)
			})

			assertTestSnapshot(t, h.Snapshot(), prev)
			if got := h.LastError(); got != prevErr {
				t.Fatalf("LastError() = %v, want preserved %v", got, prevErr)
			}

			change, err := h.Apply(testConfig{Name: "recovered", Limit: 2})
			if err != nil {
				t.Fatalf("Apply() after panic error = %v", err)
			}
			if !change.Changed {
				t.Fatal("Apply() after panic changed = false, want true")
			}
			if got, want := h.Snapshot().Value.Name, "recovered"; got != want {
				t.Fatalf("Snapshot().Value.Name = %q, want %q", got, want)
			}
			if got := h.LastError(); got != nil {
				t.Fatalf("LastError() after recovery = %v, want nil", got)
			}
		})
	}
}

func TestConcurrentPublishedApplyAdvancesOncePerSuccess(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial"})
	const writers = 32

	runConcurrentErrors(t, writers, func(i int) error {
		_, err := h.Apply(testConfig{Name: "published", Limit: i + 1})
		return err
	})

	if got, want := h.Revision(), snapshot.Revision(1+writers); got != want {
		t.Fatalf("Revision() = %d, want %d", got, want)
	}
}

func TestConcurrentRejectedApplyDoesNotAdvanceRevision(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1})
	prev := h.Revision()

	runConcurrent(32, func(int) {
		_, _ = h.Apply(testConfig{Name: "bad", Limit: -1})
	})

	if got := h.Revision(); got != prev {
		t.Fatalf("Revision() = %d, want %d", got, prev)
	}
	if got, want := h.Snapshot().Value.Name, "initial"; got != want {
		t.Fatalf("Snapshot().Value.Name = %q, want %q", got, want)
	}
}

func TestConcurrentInvalidApplyKeepsLastGood(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1})
	prev := h.Snapshot()

	runConcurrent(32, func(int) {
		_, _ = h.Apply(testConfig{Name: "bad", Limit: -1})
	})

	cur := h.Snapshot()
	if cur.Revision != prev.Revision {
		t.Fatalf("Revision() = %d, want %d", cur.Revision, prev.Revision)
	}
	if got, want := cur.Value.Name, "initial"; got != want {
		t.Fatalf("current name = %q, want %q", got, want)
	}
	if h.LastError() == nil {
		t.Fatal("LastError() = nil, want rejected apply error")
	}
}

func TestConcurrentMixedApplyKeepsValidLastGoodSnapshot(t *testing.T) {
	var published atomic.Int64
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 0})

	runConcurrent(64, func(i int) {
		limit := i
		if i%2 == 0 {
			limit = -1
		}
		change, err := h.Apply(testConfig{Name: "candidate", Limit: limit})
		if err == nil && change.Changed {
			published.Add(1)
		}
	})

	snap := h.Snapshot()
	if snap.IsZeroRevision() {
		t.Fatal("Snapshot().Revision is zero")
	}
	if snap.Value.Limit < 0 {
		t.Fatalf("Snapshot().Value.Limit = %d, want non-negative", snap.Value.Limit)
	}
	if got, want := h.Revision(), snapshot.Revision(1+published.Load()); got != want {
		t.Fatalf("Revision() = %d, want %d", got, want)
	}
}

func BenchmarkHolderApplyPublished(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := range b.N {
		benchmarkChangeSink, benchmarkErrorSink = h.Apply(testConfig{Name: "next", Limit: i})
	}
}

func BenchmarkHolderApplyEqualNoop(b *testing.B) {
	h := newBenchmarkHolder(b, WithEqual(equalTestConfig))

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		benchmarkChangeSink, benchmarkErrorSink = h.Apply(testConfig{Name: "initial", Limit: 1})
	}
}

func BenchmarkHolderApplyEqualChanged(b *testing.B) {
	h := newBenchmarkHolder(b, WithEqual(equalTestConfig))

	b.ReportAllocs()
	b.ResetTimer()

	for i := range b.N {
		benchmarkChangeSink, benchmarkErrorSink = h.Apply(testConfig{Name: "next", Limit: i + 2})
	}
}

func BenchmarkHolderApplyValidateFailed(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		benchmarkChangeSink, benchmarkErrorSink = h.Apply(testConfig{Name: "bad", Limit: -1})
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
		benchmarkChangeSink, benchmarkErrorSink = h.Apply(testConfig{Name: "bad", Limit: 1})
	}
}

func BenchmarkHolderApplyWithNormalizerSuccess(b *testing.B) {
	h := newBenchmarkHolder(b, WithNormalizer(func(cfg testConfig) (testConfig, error) {
		cfg.Name = "normalized"
		cfg.Limit++
		return cfg, nil
	}))

	b.ReportAllocs()
	b.ResetTimer()

	for i := range b.N {
		benchmarkChangeSink, benchmarkErrorSink = h.Apply(testConfig{Name: "raw", Limit: i})
	}
}

func BenchmarkHolderApplyCloneSlice100(b *testing.B) {
	h := newBenchmarkHolder(b)
	candidate := benchmarkConfigWithTags(100)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		benchmarkChangeSink, benchmarkErrorSink = h.Apply(candidate)
	}
}

func BenchmarkHolderApplyPublishedParallel(b *testing.B) {
	h := newBenchmarkHolder(b)
	var next atomic.Int64

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var local changeLocal
		for pb.Next() {
			limit := int(next.Add(1))
			local.change, local.err = h.Apply(testConfig{Name: "next", Limit: limit})
		}
		local.keep()
	})
}

func BenchmarkHolderApplyRejectedParallel(b *testing.B) {
	h := newBenchmarkHolder(b)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var local changeLocal
		for pb.Next() {
			local.change, local.err = h.Apply(testConfig{Name: "bad", Limit: -1})
		}
		local.keep()
	})
}

func BenchmarkHolderApplyEqualNoopParallel(b *testing.B) {
	h := newBenchmarkHolder(b, WithEqual(equalTestConfig))

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var local changeLocal
		for pb.Next() {
			local.change, local.err = h.Apply(testConfig{Name: "initial", Limit: 1})
		}
		local.keep()
	})
}
