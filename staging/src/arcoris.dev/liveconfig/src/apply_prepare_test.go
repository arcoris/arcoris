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
	"slices"
	"testing"
)

func TestPreparePipelineOrderCloneNormalizeValidate(t *testing.T) {
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
	h, err := New(testConfig{Name: "initial"}, WithClone(clone), WithNormalizer(normalize), WithValidator(validate))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	calls = nil
	candidate, reason, err := h.prepare(testConfig{Name: "next"})
	if err != nil {
		t.Fatalf("prepare() error = %v", err)
	}
	if reason != ChangeReasonUnknown {
		t.Fatalf("prepare() reason = %s, want %s", reason, ChangeReasonUnknown)
	}
	if want := []string{"clone", "normalize", "validate"}; !slices.Equal(calls, want) {
		t.Fatalf("calls = %#v, want %#v", calls, want)
	}
	if got, want := candidate.Name, "next-cloned-normalized"; got != want {
		t.Fatalf("candidate name = %q, want %q", got, want)
	}
}

func TestPrepareNormalizeFailureDoesNotCallValidator(t *testing.T) {
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
	)
	calls = nil

	_, reason, err := h.prepare(testConfig{Name: "bad"})
	if err != errNormalize {
		t.Fatalf("prepare() error = %v, want %v", err, errNormalize)
	}
	if reason != ChangeReasonNormalizeFailed {
		t.Fatalf("prepare() reason = %s, want %s", reason, ChangeReasonNormalizeFailed)
	}
	if want := []string{"normalize"}; !slices.Equal(calls, want) {
		t.Fatalf("calls = %#v, want %#v", calls, want)
	}
}

func TestPrepareValidateFailureReportsValidateReason(t *testing.T) {
	errValidate := errors.New("validate failed")
	h := newTestHolder(
		t,
		testConfig{Name: "initial"},
		WithValidator(func(cfg testConfig) error {
			if cfg.Name == "bad" {
				return errValidate
			}
			return nil
		}),
	)

	_, reason, err := h.prepare(testConfig{Name: "bad"})
	if err != errValidate {
		t.Fatalf("prepare() error = %v, want %v", err, errValidate)
	}
	if reason != ChangeReasonValidateFailed {
		t.Fatalf("prepare() reason = %s, want %s", reason, ChangeReasonValidateFailed)
	}
}
