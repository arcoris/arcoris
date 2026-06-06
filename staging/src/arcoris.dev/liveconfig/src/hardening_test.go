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
	"maps"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"arcoris.dev/snapshot"
)

type mutableConfig struct {
	Name   string
	Tags   []string
	Labels map[string]string
	Nested *nestedConfig
}

type nestedConfig struct {
	Values []string
}

type mutableClock struct {
	now time.Time
}

func (c *mutableClock) Now() time.Time {
	return c.now
}

func (c *mutableClock) Since(t time.Time) time.Duration {
	return c.now.Sub(t)
}

func (c *mutableClock) Until(t time.Time) time.Duration {
	return t.Sub(c.now)
}

func cloneMutableConfig(v mutableConfig) mutableConfig {
	out := mutableConfig{
		Name:   v.Name,
		Tags:   slices.Clone(v.Tags),
		Labels: maps.Clone(v.Labels),
	}
	if v.Nested != nil {
		out.Nested = &nestedConfig{
			Values: slices.Clone(v.Nested.Values),
		}
	}
	return out
}

func equalMutableConfig(a, b mutableConfig) bool {
	return a.Name == b.Name &&
		slices.Equal(a.Tags, b.Tags) &&
		maps.Equal(a.Labels, b.Labels) &&
		equalNestedConfig(a.Nested, b.Nested)
}

func equalNestedConfig(a, b *nestedConfig) bool {
	switch {
	case a == nil || b == nil:
		return a == b
	default:
		return slices.Equal(a.Values, b.Values)
	}
}

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

func TestNewRejectsInitialNormalizeFailure(t *testing.T) {
	errNormalize := errors.New("normalize failed")
	h, err := New(testConfig{Name: "initial"}, WithNormalizer(func(testConfig) (testConfig, error) {
		return testConfig{}, errNormalize
	}))
	if err != errNormalize {
		t.Fatalf("New() error = %v, want %v", err, errNormalize)
	}
	if h != nil {
		t.Fatalf("New() holder = %#v, want nil", h)
	}
}

func TestNewRejectsInitialValidateFailureAfterNormalization(t *testing.T) {
	errValidate := errors.New("validate failed")
	validatorSawNormalized := false
	h, err := New(
		testConfig{Name: "raw"},
		WithNormalizer(func(cfg testConfig) (testConfig, error) {
			cfg.Name = "normalized"
			return cfg, nil
		}),
		WithValidator(func(cfg testConfig) error {
			validatorSawNormalized = cfg.Name == "normalized"
			return errValidate
		}),
	)
	if err != errValidate {
		t.Fatalf("New() error = %v, want %v", err, errValidate)
	}
	if h != nil {
		t.Fatalf("New() holder = %#v, want nil", h)
	}
	if !validatorSawNormalized {
		t.Fatal("validator did not see normalized initial config")
	}
}

func TestNewPublishesNormalizedInitialConfig(t *testing.T) {
	h, err := New(testConfig{Name: "raw"}, WithNormalizer(func(cfg testConfig) (testConfig, error) {
		cfg.Name = "normalized"
		cfg.Limit = 7
		return cfg, nil
	}))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got, want := h.Snapshot().Value, (testConfig{Name: "normalized", Limit: 7}); got.Name != want.Name || got.Limit != want.Limit {
		t.Fatalf("Snapshot().Value = %#v, want %#v", got, want)
	}
}

func TestNewSuccessfulHolderStartsAtRevisionOne(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial"})
	if got, want := h.Revision(), snapshot.ZeroRevision.Next(); got != want {
		t.Fatalf("Revision() = %d, want %d", got, want)
	}
}

func TestNewSuccessfulHolderHasNilLastError(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial"})
	if got := h.LastError(); got != nil {
		t.Fatalf("LastError() = %v, want nil", got)
	}
}

func TestApplyPipelineOrderCloneNormalizeValidateEqualPublish(t *testing.T) {
	var calls []string
	clone := func(cfg testConfig) testConfig {
		calls = append(calls, "clone")
		cfg.Name += "-cloned"
		return cfg
	}
	normalize := func(cfg testConfig) (testConfig, error) {
		calls = append(calls, "normalize")
		cfg.Name += "-normalized"
		return cfg, nil
	}
	validate := func(testConfig) error {
		calls = append(calls, "validate")
		return nil
	}
	equal := func(testConfig, testConfig) bool {
		calls = append(calls, "equal")
		return false
	}
	h, err := New(testConfig{Name: "initial"}, WithClone(clone), WithNormalizer(normalize), WithValidator(validate), WithEqual(equal))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	calls = nil
	change, err := h.Apply(testConfig{Name: "next"})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if !change.Changed {
		t.Fatal("Apply().Changed = false, want true")
	}
	if want := []string{"clone", "normalize", "validate", "equal"}; !slices.Equal(calls, want) {
		t.Fatalf("calls = %#v, want %#v", calls, want)
	}
	if got, want := h.Snapshot().Value.Name, "next-cloned-normalized"; got != want {
		t.Fatalf("published name = %q, want %q", got, want)
	}
}

func TestApplyNormalizeFailureDoesNotCallValidatorOrEqual(t *testing.T) {
	errNormalize := errors.New("normalize failed")
	var calls []string
	h := newTestHolder(
		t,
		testConfig{Name: "initial"},
		WithNormalizer(func(cfg testConfig) (testConfig, error) {
			calls = append(calls, "normalize")
			if cfg.Name == "bad" {
				return testConfig{}, errNormalize
			}
			return cfg, nil
		}),
		WithValidator(func(testConfig) error {
			calls = append(calls, "validate")
			return nil
		}),
		WithEqual(func(testConfig, testConfig) bool {
			calls = append(calls, "equal")
			return false
		}),
	)
	calls = nil

	_, err := h.Apply(testConfig{Name: "bad"})
	if err != errNormalize {
		t.Fatalf("Apply() error = %v, want %v", err, errNormalize)
	}
	if want := []string{"normalize"}; !slices.Equal(calls, want) {
		t.Fatalf("calls = %#v, want %#v", calls, want)
	}
}

func TestApplyValidateFailureDoesNotCallEqual(t *testing.T) {
	errValidate := errors.New("validate failed")
	var calls []string
	h := newTestHolder(
		t,
		testConfig{Name: "initial"},
		WithValidator(func(cfg testConfig) error {
			calls = append(calls, "validate")
			if cfg.Name == "bad" {
				return errValidate
			}
			return nil
		}),
		WithEqual(func(testConfig, testConfig) bool {
			calls = append(calls, "equal")
			return false
		}),
	)
	calls = nil

	_, err := h.Apply(testConfig{Name: "bad"})
	if err != errValidate {
		t.Fatalf("Apply() error = %v, want %v", err, errValidate)
	}
	if want := []string{"validate"}; !slices.Equal(calls, want) {
		t.Fatalf("calls = %#v, want %#v", calls, want)
	}
}

func TestEqualSeesPreparedCandidateAndCurrentPublishedValue(t *testing.T) {
	var current, candidate testConfig
	h := newTestHolder(
		t,
		testConfig{Name: "current", Tags: []string{"published"}},
		WithNormalizer(func(cfg testConfig) (testConfig, error) {
			cfg.Tags = append(cfg.Tags, "prepared")
			return cfg, nil
		}),
		WithEqual(func(a, b testConfig) bool {
			current = a
			candidate = b
			return false
		}),
	)

	_, err := h.Apply(testConfig{Name: "next", Tags: []string{"candidate"}})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if got, want := current.Name, "current"; got != want {
		t.Fatalf("equal current name = %q, want %q", got, want)
	}
	if got, want := candidate.Tags, []string{"candidate", "prepared"}; !slices.Equal(got, want) {
		t.Fatalf("equal candidate tags = %#v, want %#v", got, want)
	}
}

func TestApplyEqualNoopDoesNotAdvanceRevision(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1}, WithEqual(equalTestConfig))
	prev := h.Revision()

	change, err := h.Apply(testConfig{Name: "initial", Limit: 1})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if change.Current.Revision != prev {
		t.Fatalf("Current.Revision = %d, want %d", change.Current.Revision, prev)
	}
	if got := h.Revision(); got != prev {
		t.Fatalf("Revision() = %d, want %d", got, prev)
	}
}

func TestApplyEqualNoopDoesNotChangeStampedUpdated(t *testing.T) {
	clk := &mutableClock{now: time.Unix(10, 0)}
	h, err := New(testConfig{Name: "initial"}, WithClock[testConfig](clk), WithEqual(equalTestConfig))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	prev := h.Stamped()

	clk.now = time.Unix(20, 0)
	_, err = h.Apply(testConfig{Name: "initial"})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if got := h.Stamped().Updated; !got.Equal(prev.Updated) {
		t.Fatalf("Stamped().Updated = %s, want %s", got, prev.Updated)
	}
}

func TestApplyEqualNoopReturnsPreviousAsCurrent(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1}, WithEqual(equalTestConfig))
	prev := h.Snapshot()

	change, err := h.Apply(testConfig{Name: "initial", Limit: 1})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	assertTestSnapshot(t, change.Current, prev)
}

func TestHolderCloneProtectsPublishedValueFromInputMutation(t *testing.T) {
	initial := mutableConfig{
		Name:   "initial",
		Tags:   []string{"a"},
		Labels: map[string]string{"k": "v"},
		Nested: &nestedConfig{Values: []string{"n"}},
	}
	h, err := New(initial, WithClone(cloneMutableConfig))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	initial.Tags[0] = "changed"
	initial.Labels["k"] = "changed"
	initial.Nested.Values[0] = "changed"
	assertMutableConfig(t, h.Snapshot().Value, mutableConfig{
		Name:   "initial",
		Tags:   []string{"a"},
		Labels: map[string]string{"k": "v"},
		Nested: &nestedConfig{Values: []string{"n"}},
	})

	next := mutableConfig{
		Name:   "next",
		Tags:   []string{"b"},
		Labels: map[string]string{"next": "value"},
		Nested: &nestedConfig{Values: []string{"m"}},
	}
	_, err = h.Apply(next)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	next.Tags[0] = "changed"
	next.Labels["next"] = "changed"
	next.Nested.Values[0] = "changed"
	assertMutableConfig(t, h.Snapshot().Value, mutableConfig{
		Name:   "next",
		Tags:   []string{"b"},
		Labels: map[string]string{"next": "value"},
		Nested: &nestedConfig{Values: []string{"m"}},
	})
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

func TestHolderBadCloneCanLeakMutableState(t *testing.T) {
	// Identity is caller misuse for mutable values: the holder cannot protect
	// published state when CloneFunc returns caller-owned storage.
	input := mutableConfig{Tags: []string{"a"}}
	h, err := New(input, WithClone(identityClone[mutableConfig]))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	input.Tags[0] = "changed"
	if got, want := h.Snapshot().Value.Tags[0], "changed"; got != want {
		t.Fatalf("Snapshot().Value.Tags[0] = %q, want %q", got, want)
	}
}

func TestHolderCloneProtectsMapsAndNestedMutableValues(t *testing.T) {
	h, err := New(
		mutableConfig{
			Labels: map[string]string{"initial": "value"},
			Nested: &nestedConfig{Values: []string{"initial"}},
		},
		WithClone(cloneMutableConfig),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	next := mutableConfig{
		Labels: map[string]string{"next": "value"},
		Nested: &nestedConfig{Values: []string{"next"}},
	}
	_, err = h.Apply(next)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	next.Labels["next"] = "changed"
	next.Nested.Values[0] = "changed"
	assertMutableConfig(t, h.Snapshot().Value, mutableConfig{
		Labels: map[string]string{"next": "value"},
		Nested: &nestedConfig{Values: []string{"next"}},
	})
}

func TestConcurrentEqualApplyDoesNotAdvanceRevision(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1}, WithEqual(equalTestConfig))
	prev := h.Revision()

	runConcurrent(32, func(int) {
		_, _ = h.Apply(testConfig{Name: "initial", Limit: 1})
	})

	if got := h.Revision(); got != prev {
		t.Fatalf("Revision() = %d, want %d", got, prev)
	}
}

func TestConcurrentPublishedApplyAdvancesOncePerSuccess(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial"})
	const writers = 32

	runConcurrent(writers, func(i int) {
		_, err := h.Apply(testConfig{Name: "published", Limit: i + 1})
		if err != nil {
			t.Errorf("Apply() error = %v", err)
		}
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

func TestZeroValueHolderPanics(t *testing.T) {
	tests := []struct {
		name string
		call func(h *Holder[testConfig])
	}{
		{name: "Snapshot", call: func(h *Holder[testConfig]) { _ = h.Snapshot() }},
		{name: "Stamped", call: func(h *Holder[testConfig]) { _ = h.Stamped() }},
		{name: "Revision", call: func(h *Holder[testConfig]) { _ = h.Revision() }},
		{name: "LastError", call: func(h *Holder[testConfig]) { _ = h.LastError() }},
		{name: "Apply", call: func(h *Holder[testConfig]) { _, _ = h.Apply(testConfig{}) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("panic = nil, want panic")
				}
			}()

			var h Holder[testConfig]
			tt.call(&h)
		})
	}
}

func runConcurrent(n int, fn func(int)) {
	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			<-start
			fn(i)
		}()
	}
	close(start)
	wg.Wait()
}

func assertTestSnapshot(t *testing.T, got, want snapshot.Snapshot[testConfig]) {
	t.Helper()
	if got.Revision != want.Revision || !equalTestConfig(got.Value, want.Value) {
		t.Fatalf("Snapshot() = %#v, want %#v", got, want)
	}
}

func assertMutableConfig(t *testing.T, got, want mutableConfig) {
	t.Helper()
	if !equalMutableConfig(got, want) {
		t.Fatalf("mutable config = %#v, want %#v", got, want)
	}
}
