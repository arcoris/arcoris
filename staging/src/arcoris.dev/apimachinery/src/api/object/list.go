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

package object

import "arcoris.dev/apimachinery/api/meta"

// List is a generic ARCORIS API object list envelope.
//
// T is the item type. It may be Object[D, O], a concrete resource type, or
// another API object representation. The list envelope validates only metadata.
type List[T any] struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty,omitzero"`

	// Items are the list payload values. api/object does not validate them.
	Items []T `json:"items"`
}

// NewList constructs a list envelope.
//
// Metadata is copied using metadata clone semantics. Items are shallow-copied
// into a new slice, but item values themselves are not cloned. The generic list
// envelope preserves nil-vs-empty item slice shape; serving layers may normalize
// API responses later if they require empty lists to encode as [].
func NewList[T any](
	typeMeta meta.TypeMeta,
	listMeta meta.ListMeta,
	items []T,
) List[T] {
	return List[T]{
		TypeMeta: typeMeta.Clone(),
		ListMeta: listMeta.Clone(),
		Items:    copyItems(items),
	}
}

// Len reports the number of list items without inspecting item values.
func (l List[T]) Len() int {
	return len(l.Items)
}

// IsZero reports whether the list envelope has no metadata and no items.
//
// It checks only slice length, not generic item contents.
func (l List[T]) IsZero() bool {
	return l.TypeMeta.IsZero() &&
		l.ListMeta.IsZero() &&
		len(l.Items) == 0
}
