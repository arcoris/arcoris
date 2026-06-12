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

package objectownership

import "arcoris.dev/apimachinery/api/fieldownership"

// Entry stores one owner and the canonical paths it owns in a surface.
//
// Duplicate entries for the same owner are valid raw Document input. Normalize
// merges them through fieldownership.State while preserving shared ownership
// across different owners.
type Entry struct {
	// Owner identifies the field owner for the listed paths.
	Owner fieldownership.Owner

	// Fields lists canonical semantic field paths owned by Owner. Parent and
	// child paths are preserved as distinct explicit ownership records.
	Fields []Path
}

// Clone returns a detached copy of e without validating or normalizing it.
//
// Clone preserves duplicate fields, raw field order, and nil-vs-empty Fields
// shape. Owner is copied by value.
func (e Entry) Clone() Entry {
	return Entry{
		Owner:  e.Owner,
		Fields: e.FieldsCopy(),
	}
}

// FieldsCopy returns a detached copy of Fields.
//
// FieldsCopy preserves nil-vs-empty slice shape and raw field order. It does
// not validate, deduplicate, sort, or canonicalize field path text.
func (e Entry) FieldsCopy() []Path {
	if e.Fields == nil {
		return nil
	}

	out := make([]Path, len(e.Fields))
	copy(out, e.Fields)

	return out
}

// IsEmpty reports whether the entry owns no fields.
//
// IsEmpty is not a validity check. It ignores Owner validity and path validity,
// and only answers whether the entry mentions fields.
//
// Empty entries can appear in raw documents, but Normalize removes them from
// canonical output and StateFromDocument ignores them through normalization.
func (e Entry) IsEmpty() bool {
	return len(e.Fields) == 0
}
