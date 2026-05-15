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

func TestValidatorAcceptsCandidate(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1})

	if _, err := h.Apply(testConfig{Name: "next", Limit: 0}); err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
}

func TestValidatorRejectsCandidate(t *testing.T) {
	errInvalid := errors.New("invalid")
	h, err := New(
		testConfig{Name: "initial"},
		WithValidator(func(cfg testConfig) error {
			if cfg.Name == "bad" {
				return errInvalid
			}
			return nil
		}),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	prev := h.Snapshot()

	change, err := h.Apply(testConfig{Name: "bad"})
	if err != errInvalid {
		t.Fatalf("Apply() error = %v, want %v", err, errInvalid)
	}
	if change.Changed {
		t.Fatal("Apply().Changed = true, want false")
	}
	if change.Current.Revision != prev.Revision {
		t.Fatalf("Current revision = %d, want %d", change.Current.Revision, prev.Revision)
	}
	if got, want := h.Snapshot().Value.Name, "initial"; got != want {
		t.Fatalf("current name = %q, want %q", got, want)
	}
	if got := h.LastError(); got != errInvalid {
		t.Fatalf("LastError() = %v, want %v", got, errInvalid)
	}
}

func TestValidatorSeesNormalizedValue(t *testing.T) {
	var seen testConfig
	normalize := func(cfg testConfig) (testConfig, error) {
		cfg.Name = "normalized"
		return cfg, nil
	}
	validate := func(cfg testConfig) error {
		seen = cfg
		if cfg.Name != "normalized" {
			return errors.New("validator saw unnormalized value")
		}
		return nil
	}

	h, err := New(testConfig{Name: "raw"}, WithNormalizer(normalize), WithValidator(validate))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got, want := seen.Name, "normalized"; got != want {
		t.Fatalf("validator saw name = %q, want %q", got, want)
	}
	if got, want := h.Snapshot().Value.Name, "normalized"; got != want {
		t.Fatalf("published name = %q, want %q", got, want)
	}
}

func TestValidatorRejectsInvalidInitialConfig(t *testing.T) {
	errInvalid := errors.New("invalid initial")

	h, err := New(
		testConfig{Name: "bad"},
		WithValidator(func(testConfig) error { return errInvalid }),
	)
	if err != errInvalid {
		t.Fatalf("New() error = %v, want %v", err, errInvalid)
	}
	if h != nil {
		t.Fatalf("New() holder = %#v, want nil", h)
	}
}
