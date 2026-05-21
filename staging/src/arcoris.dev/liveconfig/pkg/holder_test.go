/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package liveconfig

import (
	"errors"
	"testing"
	"time"

	"arcoris.dev/snapshot"
)

type testConfig struct {
	Name  string
	Limit int
	Tags  []string
}

type testClock struct {
	now time.Time
}

func (c testClock) Now() time.Time {
	return c.now
}

func (c testClock) Since(t time.Time) time.Duration {
	return c.now.Sub(t)
}

func newTestClock() testClock {
	return testClock{now: time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)}
}

func cloneTestConfig(cfg testConfig) testConfig {
	cfg.Tags = append([]string(nil), cfg.Tags...)
	return cfg
}

func validTestConfig(cfg testConfig) error {
	if cfg.Limit < 0 {
		return errors.New("limit must be non-negative")
	}
	return nil
}

func equalTestConfig(a, b testConfig) bool {
	if a.Name != b.Name || a.Limit != b.Limit || len(a.Tags) != len(b.Tags) {
		return false
	}
	for i := range a.Tags {
		if a.Tags[i] != b.Tags[i] {
			return false
		}
	}
	return true
}

func newTestHolder(t *testing.T, initial testConfig, opts ...Option[testConfig]) *Holder[testConfig] {
	t.Helper()

	base := []Option[testConfig]{
		WithClock[testConfig](newTestClock()),
		WithClone(cloneTestConfig),
		WithValidator(validTestConfig),
	}
	base = append(base, opts...)

	h, err := New(initial, base...)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return h
}

func TestNewPublishesInitialConfig(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 3, Tags: []string{"a"}})

	snap := h.Snapshot()
	if snap.IsZeroRevision() {
		t.Fatal("Snapshot revision is zero")
	}
	if got, want := snap.Revision, snapshot.ZeroRevision.Next(); got != want {
		t.Fatalf("Snapshot revision = %d, want %d", got, want)
	}
	if got, want := snap.Value.Name, "initial"; got != want {
		t.Fatalf("Snapshot value name = %q, want %q", got, want)
	}
	if got, want := snap.Value.Limit, 3; got != want {
		t.Fatalf("Snapshot value limit = %d, want %d", got, want)
	}
	if err := h.LastError(); err != nil {
		t.Fatalf("LastError() = %v, want nil", err)
	}
}

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

func TestNewRejectsInvalidInitialConfig(t *testing.T) {
	h, err := New(
		testConfig{Name: "invalid", Limit: -1},
		WithValidator(validTestConfig),
	)
	if err == nil {
		t.Fatal("New() error = nil, want validation error")
	}
	if h != nil {
		t.Fatalf("New() holder = %#v, want nil", h)
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

func TestHolderImplementsSnapshotContracts(t *testing.T) {
	var _ snapshot.Source[testConfig] = (*Holder[testConfig])(nil)
	var _ snapshot.StampedSource[testConfig] = (*Holder[testConfig])(nil)
	var _ snapshot.RevisionSource = (*Holder[testConfig])(nil)
}

func TestNilHolderPanics(t *testing.T) {
	tests := []struct {
		name string
		call func()
	}{
		{
			name: "Snapshot",
			call: func() {
				var h *Holder[testConfig]
				_ = h.Snapshot()
			},
		},
		{
			name: "Stamped",
			call: func() {
				var h *Holder[testConfig]
				_ = h.Stamped()
			},
		},
		{
			name: "Revision",
			call: func() {
				var h *Holder[testConfig]
				_ = h.Revision()
			},
		},
		{
			name: "LastError",
			call: func() {
				var h *Holder[testConfig]
				_ = h.LastError()
			},
		},
		{
			name: "Apply",
			call: func() {
				var h *Holder[testConfig]
				_, _ = h.Apply(testConfig{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if got := recover(); got != ErrNilHolder {
					t.Fatalf("panic = %v, want %v", got, ErrNilHolder)
				}
			}()

			tt.call()
		})
	}
}
