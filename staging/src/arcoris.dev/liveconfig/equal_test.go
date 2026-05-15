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

func TestEqualValueWithoutEqualFuncReportsFalse(t *testing.T) {
	cfg := defaultConfig[testConfig]()
	if equalValue(cfg, testConfig{Name: "a"}, testConfig{Name: "a"}) {
		t.Fatal("equalValue() = true, want false without EqualFunc")
	}
}

func TestEqualValueUsesEqualFunc(t *testing.T) {
	cfg := newConfig(WithEqual(func(a, b testConfig) bool { return a.Name == b.Name }))
	if !equalValue(cfg, testConfig{Name: "a"}, testConfig{Name: "a"}) {
		t.Fatal("equalValue() = false, want true")
	}
	if equalValue(cfg, testConfig{Name: "a"}, testConfig{Name: "b"}) {
		t.Fatal("equalValue() = true, want false")
	}
}

func TestEqualFuncTrueSuppressesPublication(t *testing.T) {
	h := newTestHolder(
		t,
		testConfig{Name: "initial", Limit: 1},
		WithEqual(func(testConfig, testConfig) bool { return true }),
	)
	prev := h.Snapshot()

	change, err := h.Apply(testConfig{Name: "next", Limit: 2})
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

func TestEqualFuncFalsePublishesCandidate(t *testing.T) {
	h := newTestHolder(
		t,
		testConfig{Name: "initial", Limit: 1},
		WithEqual(func(testConfig, testConfig) bool { return false }),
	)
	prev := h.Snapshot()

	change, err := h.Apply(testConfig{Name: "initial", Limit: 1})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if !change.Changed {
		t.Fatal("Apply().Changed = false, want true")
	}
	if got, want := change.Current.Revision, prev.Revision.Next(); got != want {
		t.Fatalf("Current revision = %d, want %d", got, want)
	}
}

func TestNilEqualAlwaysPublishesValidCandidate(t *testing.T) {
	h := newTestHolder(t, testConfig{Name: "initial", Limit: 1})
	prev := h.Snapshot()

	change, err := h.Apply(testConfig{Name: "initial", Limit: 1})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if !change.Changed {
		t.Fatal("Apply().Changed = false, want true")
	}
	if got, want := change.Current.Revision, prev.Revision.Next(); got != want {
		t.Fatalf("Current revision = %d, want %d", got, want)
	}
}

func TestEqualEvaluatesAfterCandidatePipeline(t *testing.T) {
	clone := func(cfg testConfig) testConfig {
		cfg.Tags = append([]string(nil), cfg.Tags...)
		cfg.Tags = append(cfg.Tags, "cloned")
		return cfg
	}
	normalize := func(cfg testConfig) (testConfig, error) {
		cfg.Tags = append(cfg.Tags, "normalized")
		return cfg, nil
	}

	validated := false
	validate := func(cfg testConfig) error {
		if len(cfg.Tags) != 2 || cfg.Tags[0] != "cloned" || cfg.Tags[1] != "normalized" {
			return errors.New("candidate was not prepared before validation")
		}
		validated = true
		return nil
	}

	equalBeforeValidator := false
	var equalCandidate testConfig
	equal := func(a, b testConfig) bool {
		equalBeforeValidator = !validated
		equalCandidate = b
		return equalTestConfig(a, b)
	}

	h, err := New(
		testConfig{Name: "initial"},
		WithClone(clone),
		WithNormalizer(normalize),
		WithValidator(validate),
		WithEqual(equal),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	validated = false
	change, err := h.Apply(testConfig{Name: "initial"})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if change.Changed {
		t.Fatal("Apply().Changed = true, want false")
	}
	if equalBeforeValidator {
		t.Fatal("equal ran before validator")
	}
	if len(equalCandidate.Tags) != 2 || equalCandidate.Tags[0] != "cloned" || equalCandidate.Tags[1] != "normalized" {
		t.Fatalf("equal candidate tags = %#v, want cloned normalized", equalCandidate.Tags)
	}
}
