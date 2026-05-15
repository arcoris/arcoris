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
)

func TestWithClockUsesConfiguredClock(t *testing.T) {
	clk := newTestClock()
	h, err := New(testConfig{Name: "initial"}, WithClock[testConfig](clk))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if got := h.Stamped().Updated; !got.Equal(clk.now) {
		t.Fatalf("Stamped().Updated = %s, want %s", got, clk.now)
	}
}

func TestWithCloneUsesConfiguredClone(t *testing.T) {
	cloneCalled := false
	clone := func(cfg testConfig) testConfig {
		cloneCalled = true
		cfg.Name = "cloned"
		return cfg
	}

	h, err := New(testConfig{Name: "initial"}, WithClone(clone))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if !cloneCalled {
		t.Fatal("clone was not called")
	}
	if got, want := h.Snapshot().Value.Name, "cloned"; got != want {
		t.Fatalf("Snapshot().Value.Name = %q, want %q", got, want)
	}
}

func TestWithNormalizerUsesConfiguredNormalizer(t *testing.T) {
	normalize := func(cfg testConfig) (testConfig, error) {
		cfg.Limit = 7
		return cfg, nil
	}

	h, err := New(testConfig{Name: "initial"}, WithNormalizer(normalize))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if got, want := h.Snapshot().Value.Limit, 7; got != want {
		t.Fatalf("Snapshot().Value.Limit = %d, want %d", got, want)
	}
}

func TestWithValidatorUsesConfiguredValidator(t *testing.T) {
	errInvalid := errors.New("invalid")

	h, err := New(
		testConfig{Name: "initial"},
		WithValidator(func(testConfig) error { return errInvalid }),
	)
	if err != errInvalid {
		t.Fatalf("New() error = %v, want %v", err, errInvalid)
	}
	if h != nil {
		t.Fatalf("New() holder = %#v, want nil", h)
	}
}

func TestWithEqualUsesConfiguredEqual(t *testing.T) {
	h := newTestHolder(
		t,
		testConfig{Name: "initial", Limit: 1},
		WithEqual(func(a, b testConfig) bool { return a.Name == b.Name }),
	)
	prev := h.Snapshot()

	change, err := h.Apply(testConfig{Name: "initial", Limit: 2})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if change.Changed {
		t.Fatal("Apply().Changed = true, want false")
	}
	if got := h.Revision(); got != prev.Revision {
		t.Fatalf("Revision() = %d, want %d", got, prev.Revision)
	}
}

func TestNilOptionsPanic(t *testing.T) {
	tests := []struct {
		name string
		call func()
		want any
	}{
		{name: "nil option", call: func() { _ = newConfig[testConfig](nil) }, want: ErrNilOption},
		{name: "nil clock", call: func() { _ = WithClock[testConfig](nil) }, want: ErrNilClock},
		{name: "nil clone", call: func() { _ = WithClone[testConfig](nil) }, want: ErrNilClone},
		{name: "nil normalizer", call: func() { _ = WithNormalizer[testConfig](nil) }, want: ErrNilNormalizer},
		{name: "nil validator", call: func() { _ = WithValidator[testConfig](nil) }, want: ErrNilValidator},
		{name: "nil equal", call: func() { _ = WithEqual[testConfig](nil) }, want: ErrNilEqual},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if got := recover(); got != tt.want {
					t.Fatalf("panic = %v, want %v", got, tt.want)
				}
			}()

			tt.call()
		})
	}
}

func TestNilOptionErrorsAreDistinct(t *testing.T) {
	errs := []error{
		ErrNilHolder,
		ErrNilOption,
		ErrNilClock,
		ErrNilClone,
		ErrNilNormalizer,
		ErrNilValidator,
		ErrNilEqual,
	}

	seen := make(map[string]struct{}, len(errs))
	for _, err := range errs {
		if err == nil {
			t.Fatal("package error is nil")
		}
		msg := err.Error()
		if _, ok := seen[msg]; ok {
			t.Fatalf("duplicate error message %q", msg)
		}
		seen[msg] = struct{}{}
	}
}
