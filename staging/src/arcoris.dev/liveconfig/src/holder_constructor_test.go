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
	"testing"

	"arcoris.dev/snapshot"
)

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
