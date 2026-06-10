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

// SelectorEntry is one field/value pair inside an associative-list selector.
//
// It does not describe an arbitrary predicate. Higher layers use entries only
// for stable associative-list identity, such as {"type":"Ready"}.
type SelectorEntry struct {
	field FieldName
	value Literal
}

// Field returns the selector field name.
func (e SelectorEntry) Field() FieldName {
	return e.field
}

// Value returns the selector literal value.
func (e SelectorEntry) Value() Literal {
	return e.value
}
