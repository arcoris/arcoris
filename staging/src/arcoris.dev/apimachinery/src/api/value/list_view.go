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

package value

// ListView exposes read-only list payload data.
//
// A view references private payload data owned by an immutable Value. Read
// methods return nested Values by value; explicit Clone methods deep-copy nested
// payloads when a caller needs a detached ownership boundary.
type ListView struct {
	// payload is the private list payload the view reads from.
	payload listPayload
}

// AsList returns a read-only list view when v is KindList.
//
// For every other kind, AsList returns ok=false. The returned view preserves item
// order without eagerly deep-cloning the whole list payload.
func (v Value) AsList() (ListView, bool) {
	if v.kind != KindList {
		return ListView{}, false
	}

	return ListView{payload: v.listValue}, true
}

// Len returns the number of list items.
func (l ListView) Len() int {
	return len(l.payload.items)
}

// IsEmpty reports whether the list has no items.
func (l ListView) IsEmpty() bool {
	return len(l.payload.items) == 0
}

// At returns an item at index without deep-cloning it.
//
// Out-of-range indexes return the zero Value and ok=false.
func (l ListView) At(index int) (Value, bool) {
	if index < 0 || index >= len(l.payload.items) {
		return Value{}, false
	}

	return l.payload.items[index], true
}

// CloneAt returns a deep-cloned item at index.
func (l ListView) CloneAt(index int) (Value, bool) {
	item, ok := l.At(index)
	if !ok {
		return Value{}, false
	}

	return item.Clone(), true
}

// Items returns a detached item slice in original order.
//
// The slice is caller-owned. Nested Values are returned by value without deep
// cloning; use CloneItems when a recursive copy is required.
func (l ListView) Items() []Value {
	if len(l.payload.items) == 0 {
		return []Value{}
	}

	items := make([]Value, len(l.payload.items))
	copy(items, l.payload.items)
	return items
}

// CloneItems returns detached list items with deep-cloned nested Values.
func (l ListView) CloneItems() []Value {
	if len(l.payload.items) == 0 {
		return []Value{}
	}

	return cloneValues(l.payload.items)
}

// ForEach visits items in original order until fn returns false.
func (l ListView) ForEach(fn func(index int, value Value) bool) {
	for i, item := range l.payload.items {
		if !fn(i, item) {
			return
		}
	}
}
