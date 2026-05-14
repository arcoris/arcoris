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

package maybe_test

import (
	"testing"

	"arcoris.dev/value/maybe"
)

func TestZeroValueIsNone(t *testing.T) {
	var m maybe.Maybe[string]

	if !m.IsNone() {
		t.Fatal("zero value must be None")
	}
	if m.IsSome() {
		t.Fatal("zero value must not be Some")
	}

	got, ok := m.Load()
	if ok {
		t.Fatalf("Load returned ok=true for zero value: %q", got)
	}
	if got != "" {
		t.Fatalf("Load returned non-zero value for None: %q", got)
	}
}

func TestSome(t *testing.T) {
	m := maybe.Some("value")

	if !m.IsSome() {
		t.Fatal("Some value must report IsSome")
	}
	if m.IsNone() {
		t.Fatal("Some value must not report IsNone")
	}

	got, ok := m.Load()
	if !ok {
		t.Fatal("Load returned ok=false for Some")
	}
	if got != "value" {
		t.Fatalf("got %q, want %q", got, "value")
	}
}

func TestNone(t *testing.T) {
	m := maybe.None[int]()

	if !m.IsNone() {
		t.Fatal("None must report IsNone")
	}
	if m.IsSome() {
		t.Fatal("None must not report IsSome")
	}

	got, ok := m.Load()
	if ok {
		t.Fatalf("Load returned ok=true for None: %d", got)
	}
	if got != 0 {
		t.Fatalf("got %d, want zero", got)
	}
}

func TestFrom(t *testing.T) {
	tests := []struct {
		name string
		ok   bool
		want bool
	}{
		{name: "some", ok: true, want: true},
		{name: "none", ok: false, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := maybe.From("value", tc.ok)
			if got := m.IsSome(); got != tc.want {
				t.Fatalf("IsSome() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFromDiscardsValueWhenAbsent(t *testing.T) {
	original := []string{"a", "b"}
	m := maybe.From(original, false)
	original[0] = "changed"

	got, ok := m.Load()
	if ok {
		t.Fatalf("Load returned ok=true for None: %#v", got)
	}
	if got != nil {
		t.Fatalf("Load returned non-zero slice for None: %#v", got)
	}
}

func TestValueOr(t *testing.T) {
	if got := maybe.Some("value").ValueOr("fallback"); got != "value" {
		t.Fatalf("Some.ValueOr() = %q, want %q", got, "value")
	}
	if got := maybe.None[string]().ValueOr("fallback"); got != "fallback" {
		t.Fatalf("None.ValueOr() = %q, want %q", got, "fallback")
	}
}

func TestMust(t *testing.T) {
	if got := maybe.Some("value").Must(); got != "value" {
		t.Fatalf("Must() = %q, want %q", got, "value")
	}
}

func TestMustPanicsForNone(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Must did not panic for None")
		}
	}()

	_ = maybe.None[string]().Must()
}
