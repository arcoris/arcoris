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

func TestNormalizerRunsBeforeValidator(t *testing.T) {
	normalize := func(cfg testConfig) (testConfig, error) {
		if cfg.Limit == 0 {
			cfg.Limit = 1
		}
		return cfg, nil
	}
	validate := func(cfg testConfig) error {
		if cfg.Limit <= 0 {
			return errors.New("limit must be positive")
		}
		return nil
	}

	h, err := New(
		testConfig{Name: "initial"},
		WithNormalizer(normalize),
		WithValidator(validate),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got, want := h.Snapshot().Value.Limit, 1; got != want {
		t.Fatalf("normalized limit = %d, want %d", got, want)
	}
}

func TestNormalizerOutputIsPublished(t *testing.T) {
	normalize := func(cfg testConfig) (testConfig, error) {
		cfg.Name = "normalized"
		cfg.Limit++
		return cfg, nil
	}
	h := newTestHolder(
		t,
		testConfig{Name: "initial", Limit: 1},
		WithNormalizer(normalize),
	)

	change, err := h.Apply(testConfig{Name: "candidate", Limit: 4})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if !change.Changed {
		t.Fatal("Apply().Changed = false, want true")
	}
	if got, want := h.Snapshot().Value.Name, "normalized"; got != want {
		t.Fatalf("published name = %q, want %q", got, want)
	}
	if got, want := h.Snapshot().Value.Limit, 5; got != want {
		t.Fatalf("published limit = %d, want %d", got, want)
	}
}

func TestNormalizerErrorRejectsCandidate(t *testing.T) {
	boom := errors.New("normalize failed")
	h := newTestHolder(
		t,
		testConfig{Name: "initial", Limit: 1},
		WithNormalizer(func(cfg testConfig) (testConfig, error) {
			if cfg.Name == "bad" {
				return testConfig{}, boom
			}
			return cfg, nil
		}),
	)
	prev := h.Snapshot()

	change, err := h.Apply(testConfig{Name: "bad", Limit: 2})
	if err != boom {
		t.Fatalf("Apply() error = %v, want %v", err, boom)
	}
	if change.Changed {
		t.Fatal("Apply().Changed = true, want false")
	}
	if change.Current.Revision != prev.Revision {
		t.Fatalf("Current revision = %d, want %d", change.Current.Revision, prev.Revision)
	}
	if h.Revision() != prev.Revision {
		t.Fatalf("Revision() = %d, want %d", h.Revision(), prev.Revision)
	}
	if got, want := h.Snapshot().Value.Name, "initial"; got != want {
		t.Fatalf("current name = %q, want %q", got, want)
	}
	if got := h.LastError(); got != boom {
		t.Fatalf("LastError() = %v, want %v", got, boom)
	}
}
