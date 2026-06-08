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

import "testing"

func TestIdentityCloneReturnsValue(t *testing.T) {
	cfg := testConfig{Name: "value", Limit: 1}
	if got := identityClone(cfg); got.Name != cfg.Name || got.Limit != cfg.Limit {
		t.Fatalf("identityClone() = %+v, want %+v", got, cfg)
	}
}

func TestCustomCloneIsUsedBeforePublication(t *testing.T) {
	clone := func(cfg testConfig) testConfig {
		cfg.Name = cfg.Name + "-clone"
		return cfg
	}

	h, err := New(testConfig{Name: "initial"}, WithClone(clone))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got, want := h.Snapshot().Value.Name, "initial-clone"; got != want {
		t.Fatalf("initial snapshot name = %q, want %q", got, want)
	}

	_, err = h.Apply(testConfig{Name: "next"})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if got, want := h.Snapshot().Value.Name, "next-clone"; got != want {
		t.Fatalf("applied snapshot name = %q, want %q", got, want)
	}
}

func TestCloneFuncProtectsPublishedValueFromInputMutation(t *testing.T) {
	input := testConfig{Name: "initial", Tags: []string{"a"}}
	h := newTestHolder(t, input)

	input.Tags[0] = "mutated"
	if got, want := h.Snapshot().Value.Tags[0], "a"; got != want {
		t.Fatalf("initial snapshot tag = %q, want %q", got, want)
	}

	next := testConfig{Name: "next", Tags: []string{"b"}}
	_, err := h.Apply(next)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	next.Tags[0] = "mutated"
	if got, want := h.Snapshot().Value.Tags[0], "b"; got != want {
		t.Fatalf("applied snapshot tag = %q, want %q", got, want)
	}
}

func TestHolderCloneProtectsPublishedMutableConfigFromInputMutation(t *testing.T) {
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
