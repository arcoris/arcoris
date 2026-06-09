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

package types

import "testing"

func TestStringViewAccessors(t *testing.T) {
	desc := String().
		MinBytes(1).
		MaxBytes(10).
		MinRunes(1).
		MaxRunes(5).
		Pattern(`^[a-z]+$`).
		Enum("red", "green").
		Descriptor()
	view := requireStringView(t, desc)

	minBytes, ok := view.MinBytes()
	requireEqual(t, ok, true)
	requireEqual(t, minBytes, 1)

	maxBytes, ok := view.MaxBytes()
	requireEqual(t, ok, true)
	requireEqual(t, maxBytes, 10)

	minRunes, ok := view.MinRunes()
	requireEqual(t, ok, true)
	requireEqual(t, minRunes, 1)

	maxRunes, ok := view.MaxRunes()
	requireEqual(t, ok, true)
	requireEqual(t, maxRunes, 5)

	pattern, ok := view.Pattern()
	requireEqual(t, ok, true)
	requireEqual(t, pattern, `^[a-z]+$`)
}

func TestStringViewEnumReturnsDetachedData(t *testing.T) {
	view := requireStringView(t, String().Enum("red", "green").Descriptor())

	enum := view.Enum()
	enum[0] = "changed"

	requireEqual(t, view.Enum()[0], "red")
}
