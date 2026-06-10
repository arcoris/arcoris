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

package fieldpath

import "testing"

func TestSelectorEntriesReturnsClone(t *testing.T) {
	selector := MustSelector(testSelectorEntry("type", StringLiteral("Ready")))
	entries := selector.Entries()

	entries[0] = testSelectorEntry("other", StringLiteral("Other"))

	got, ok := selector.Get(testField("type"))
	requireEqual(t, ok, true)
	requireEqual(t, got.Equal(StringLiteral("Ready")), true)
}

func TestSelectorEntryReturnsEntryByIndex(t *testing.T) {
	selector := MustSelector(testSelectorEntry("type", StringLiteral("Ready")))

	entry, ok := selector.Entry(0)
	requireEqual(t, ok, true)
	requireEqual(t, entry.Field(), testField("type"))

	_, ok = selector.Entry(1)
	requireEqual(t, ok, false)
}

func TestSelectorFieldsReturnsClone(t *testing.T) {
	selector := MustSelector(testSelectorEntry("type", StringLiteral("Ready")))
	fields := selector.Fields()

	fields[0] = testField("status")

	requireEqual(t, selector.Has(testField("type")), true)
	requireEqual(t, selector.Has(testField("status")), false)
}

func TestSelectorForEachStopsEarly(t *testing.T) {
	selector := MustSelector(
		testSelectorEntry("type", StringLiteral("Ready")),
		testSelectorEntry("name", StringLiteral("main")),
	)
	visited := 0

	selector.ForEach(func(index int, entry SelectorEntry) bool {
		visited++
		return false
	})

	requireEqual(t, visited, 1)
}

func TestSelectorGet(t *testing.T) {
	selector := MustSelector(testSelectorEntry("type", StringLiteral("Ready")))

	got, ok := selector.Get(testField("type"))
	requireEqual(t, ok, true)
	requireEqual(t, got.Equal(StringLiteral("Ready")), true)
}

func TestSelectorHas(t *testing.T) {
	selector := MustSelector(testSelectorEntry("type", StringLiteral("Ready")))

	requireEqual(t, selector.Has(testField("type")), true)
	requireEqual(t, selector.Has(testField("status")), false)
}
