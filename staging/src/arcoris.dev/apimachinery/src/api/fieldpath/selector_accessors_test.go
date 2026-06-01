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
	selector := MustSelector(NewSelectorEntry("type", StringLiteral("Ready")))
	entries := selector.Entries()

	entries[0] = NewSelectorEntry("other", StringLiteral("Other"))

	got, ok := selector.Get("type")
	requireEqual(t, ok, true)
	requireEqual(t, got.Equal(StringLiteral("Ready")), true)
}

func TestSelectorGet(t *testing.T) {
	selector := MustSelector(NewSelectorEntry("type", StringLiteral("Ready")))

	got, ok := selector.Get("type")
	requireEqual(t, ok, true)
	requireEqual(t, got.Equal(StringLiteral("Ready")), true)
}

func TestSelectorHas(t *testing.T) {
	selector := MustSelector(NewSelectorEntry("type", StringLiteral("Ready")))

	requireEqual(t, selector.Has("type"), true)
	requireEqual(t, selector.Has("status"), false)
}
