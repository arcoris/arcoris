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

package codecregistry

import "testing"

func TestNewEntryID(t *testing.T) {
	tests := []string{
		"json.public",
		"json-storage",
		"json_storage",
		"codec/json/public",
		"codec/json-public_1",
	}

	for _, value := range tests {
		t.Run(value, func(t *testing.T) {
			id, err := NewEntryID(value)
			requireNoError(t, err)
			if id.String() != value {
				t.Fatalf("String() = %q; want %q", id.String(), value)
			}
		})
	}
}

func TestNewEntryIDRejectsEmpty(t *testing.T) {
	_, err := NewEntryID("")

	requireErrorIs(t, err, ErrInvalidEntryID)
}

func TestNewEntryIDRejectsUppercase(t *testing.T) {
	_, err := NewEntryID("JSON.Public")

	requireErrorIs(t, err, ErrInvalidEntryID)
}

func TestNewEntryIDRejectsWhitespace(t *testing.T) {
	tests := []string{
		"json public",
		" json.public",
		"json.public ",
	}

	for _, value := range tests {
		t.Run(value, func(t *testing.T) {
			_, err := NewEntryID(value)
			requireErrorIs(t, err, ErrInvalidEntryID)
		})
	}
}

func TestNewEntryIDRejectsControlCharacters(t *testing.T) {
	_, err := NewEntryID("json\npublic")

	requireErrorIs(t, err, ErrInvalidEntryID)
}

func TestNewEntryIDRejectsInvalidSeparators(t *testing.T) {
	tests := []string{
		".json",
		"json.",
		"json..public",
		"json//public",
	}

	for _, value := range tests {
		t.Run(value, func(t *testing.T) {
			_, err := NewEntryID(value)
			requireErrorIs(t, err, ErrInvalidEntryID)
		})
	}
}

func TestEntryIDNormalize(t *testing.T) {
	id, err := EntryID("json.public").Normalize()
	requireNoError(t, err)

	if id != MustEntryID("json.public") {
		t.Fatalf("Normalize() = %q", id)
	}
}

func TestEntryIDIsZero(t *testing.T) {
	if !EntryID("").IsZero() {
		t.Fatalf("empty ID IsZero() = false")
	}
	if MustEntryID("json.public").IsZero() {
		t.Fatalf("non-empty ID IsZero() = true")
	}
}

func TestMustEntryIDPanicsOnInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("MustEntryID() did not panic")
		}
	}()

	_ = MustEntryID("JSON.Public")
}
