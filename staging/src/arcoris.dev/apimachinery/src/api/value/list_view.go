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
// A view owns a cloned payload snapshot. Methods that return items clone again
// so callers cannot mutate the source Value or the view through returned data.
type ListView struct {
	// payload is a detached list payload snapshot.
	payload listPayload
}

// List returns a detached list view when v is KindList.
//
// For every other kind, List returns ok=false. The returned view preserves item
// order.
func (v Value) List() (ListView, bool) {
	if v.kind != KindList {
		return ListView{}, false
	}

	return ListView{payload: cloneListPayload(v.listValue)}, true
}

// Len returns the number of list items.
func (l ListView) Len() int {
	return len(l.payload.items)
}

// IsEmpty reports whether the list has no items.
func (l ListView) IsEmpty() bool {
	return len(l.payload.items) == 0
}

// At returns a cloned item at index.
//
// Out-of-range indexes return the zero Value and ok=false. Returned items are
// detached clones.
func (l ListView) At(index int) (Value, bool) {
	if index < 0 || index >= len(l.payload.items) {
		return Value{}, false
	}

	return l.payload.items[index].Clone(), true
}

// Items returns detached list items in original order.
//
// The slice and every nested Value are cloned, preserving immutable-by-
// convention behavior for list views.
func (l ListView) Items() []Value {
	return cloneValues(l.payload.items)
}
