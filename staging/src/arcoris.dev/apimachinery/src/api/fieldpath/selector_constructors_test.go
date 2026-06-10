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

func TestNewSelectorSortsEntriesByField(t *testing.T) {
	selector, err := NewSelector(
		testSelectorEntry("port", Uint64Literal(443)),
		testSelectorEntry("host", StringLiteral("api.example.com")),
	)
	requireNoError(t, err)

	entries := selector.Entries()
	requireEqual(t, entries[0].Field(), testField("host"))
	requireEqual(t, entries[1].Field(), testField("port"))
}

func TestNewSelectorClonesEntries(t *testing.T) {
	entries := []SelectorEntry{testSelectorEntry("type", StringLiteral("Ready"))}

	selector, err := NewSelector(entries...)
	requireNoError(t, err)

	entries[0] = testSelectorEntry("type", StringLiteral("Changed"))

	got, ok := selector.Get(testField("type"))
	requireEqual(t, ok, true)
	requireEqual(t, got.Equal(StringLiteral("Ready")), true)
}

func TestNewSelectorRejectsEmptySelector(t *testing.T) {
	_, err := NewSelector()

	requireErrorIs(t, err, ErrInvalidSelector)
	requireErrorIs(t, err, ErrEmptySelector)
}

func TestNewSelectorRejectsEmptyEntryField(t *testing.T) {
	_, err := NewSelector(NewSelectorEntry(FieldName(""), StringLiteral("Ready")))

	requireErrorIs(t, err, ErrInvalidSelector)
	requireErrorIs(t, err, ErrInvalidEntry)
	requireErrorIs(t, err, ErrEmptyFieldName)
}

func TestNewSelectorRejectsDuplicateEntryField(t *testing.T) {
	_, err := NewSelector(
		testSelectorEntry("type", StringLiteral("Ready")),
		testSelectorEntry("type", StringLiteral("Scheduled")),
	)

	requireErrorIs(t, err, ErrInvalidSelector)
	requireErrorIs(t, err, ErrDuplicateSelectorField)
}

func TestMustSelectorPanicsOnInvalidSelector(t *testing.T) {
	requirePanic(t, func() {
		MustSelector()
	})
}
