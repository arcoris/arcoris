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
