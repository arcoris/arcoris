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

// MapView exposes read-only map payload data.
//
// A view owns a cloned payload snapshot. Methods that return values or slices
// clone again, so callers cannot mutate the original Value or the view through
// returned data.
type MapView struct {
	// payload is a detached map payload snapshot.
	payload mapPayload
}

// Map returns a detached map view when v is KindMap.
//
// For every other kind, Map returns ok=false. The returned view is detached from
// v and preserves entry order.
func (v Value) Map() (MapView, bool) {
	if v.kind != KindMap {
		return MapView{}, false
	}

	return MapView{payload: cloneMapPayload(v.mapValue)}, true
}

// Len returns the number of map entries.
func (m MapView) Len() int {
	return len(m.payload.entries)
}

// IsEmpty reports whether the map has no entries.
func (m MapView) IsEmpty() bool {
	return len(m.payload.entries) == 0
}

// Has reports whether the map contains key.
//
// Lookup is linear by design. The value model favors cheap construction and
// cloning over storing per-map indexes for small payloads.
func (m MapView) Has(key string) bool {
	return findMapEntry(m.payload.entries, key) >= 0
}

// Get returns a cloned entry value by key.
//
// The returned Value is detached from the view. Missing keys return the zero
// Value and ok=false.
func (m MapView) Get(key string) (Value, bool) {
	index := findMapEntry(m.payload.entries, key)
	if index < 0 {
		return Value{}, false
	}

	return m.payload.entries[index].Value.Clone(), true
}

// Entries returns detached entries in original order.
//
// Both the slice and every nested Value are cloned. Mutating the result cannot
// affect the view or source Value.
func (m MapView) Entries() []Entry {
	return cloneEntries(m.payload.entries)
}

// Keys returns detached keys in original order.
//
// Keys mirrors Entries order and returns a caller-owned slice.
func (m MapView) Keys() []string {
	keys := make([]string, 0, len(m.payload.entries))
	for _, entry := range m.payload.entries {
		keys = append(keys, entry.Key)
	}

	return keys
}

// findMapEntry returns the ordered entry index for key.
//
// A negative result means the key is absent. Keeping the search as a tiny helper
// makes Has and Get share the same linear lookup semantics.
func findMapEntry(entries []Entry, key string) int {
	for i, entry := range entries {
		if entry.Key == key {
			return i
		}
	}

	return -1
}
