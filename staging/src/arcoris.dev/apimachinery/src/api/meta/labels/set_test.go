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

package labels

import (
	"slices"
	"testing"
)

func TestSet(t *testing.T) {
	var nilSet Set
	if !nilSet.IsZero() {
		t.Fatal("nil Set IsZero() = false")
	}
	if nilSet.Len() != 0 {
		t.Fatalf("nil Set Len() = %d, want 0", nilSet.Len())
	}

	set := Set{"role": "worker"}
	if set.IsZero() {
		t.Fatal("non-zero Set IsZero() = true")
	}
	if set.Len() != 1 {
		t.Fatalf("Len() = %d, want 1", set.Len())
	}
	if !set.Has("role") {
		t.Fatal("Has() = false")
	}
	if value, ok := set.Get("role"); !ok || value != "worker" {
		t.Fatalf("Get() = %q, %v", value, ok)
	}
}

func TestSetUsesTypedKeysAndValues(t *testing.T) {
	var set Set = map[Key]Value{
		Key("role"): Value("worker"),
	}

	value, ok := set.Get(Key("role"))
	if !ok || value != Value("worker") {
		t.Fatalf("Get() = %q, %v", value, ok)
	}
}

func TestSetKeysReturnsSortedKeys(t *testing.T) {
	set := Set{
		"tier": "backend",
		"role": "worker",
		"app":  "scheduler",
	}

	got := set.Keys()
	want := []Key{"app", "role", "tier"}
	if !slices.Equal(got, want) {
		t.Fatalf("Keys() = %#v; want %#v", got, want)
	}
}

func TestSetStringConversionsDetachStorage(t *testing.T) {
	raw := map[string]string{"role": "worker"}
	set, err := FromStrings(raw)
	requireNoError(t, err)

	raw["role"] = "manager"
	if set["role"] != "worker" {
		t.Fatal("FromStrings result aliases input map")
	}

	strings := set.Strings()
	strings["role"] = "manager"
	if set["role"] != "worker" {
		t.Fatal("Strings result aliases typed set")
	}
}

func TestSetFromStringsValidatesValues(t *testing.T) {
	_, err := FromStrings(map[string]string{"Role": "worker"})
	requireErrorIs(t, err, ErrInvalidSet)

	_, err = FromStrings(map[string]string{"role": "worker value"})
	requireErrorIs(t, err, ErrInvalidSet)
}
