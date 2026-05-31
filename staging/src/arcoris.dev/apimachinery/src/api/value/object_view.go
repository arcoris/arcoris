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

// ObjectView exposes read-only object payload data.
//
// A view references private immutable-by-convention payload data. Methods that
// return values or slices clone those results, so callers cannot mutate the
// original Value or the view through returned data.
type ObjectView struct {
	// payload is the private object payload the view reads from.
	payload objectPayload
}

// Object returns a read-only object view when v is KindObject.
//
// For every other kind, Object returns ok=false. The returned view preserves
// object field order without eagerly deep-cloning the whole object payload.
func (v Value) Object() (ObjectView, bool) {
	if v.kind != KindObject {
		return ObjectView{}, false
	}

	return ObjectView{payload: v.objectValue}, true
}

// Len returns the number of object fields.
func (o ObjectView) Len() int {
	return len(o.payload.fields)
}

// IsEmpty reports whether the object has no fields.
func (o ObjectView) IsEmpty() bool {
	return len(o.payload.fields) == 0
}

// Has reports whether the object contains name.
//
// Lookup is linear by design. The value model favors cheap construction and
// cloning over storing per-object indexes for small payloads.
func (o ObjectView) Has(name string) bool {
	return findObjectField(o.payload.fields, name) >= 0
}

// Get returns a cloned field value by name.
//
// The returned Value is detached from the view. Missing names return the zero
// Value and ok=false.
func (o ObjectView) Get(name string) (Value, bool) {
	index := findObjectField(o.payload.fields, name)
	if index < 0 {
		return Value{}, false
	}

	return o.payload.fields[index].Value.Clone(), true
}

// Fields returns detached fields in original order.
//
// Both the slice and every nested Value are cloned. Mutating the result cannot
// affect the view or source Value.
func (o ObjectView) Fields() []Field {
	if len(o.payload.fields) == 0 {
		return []Field{}
	}

	return cloneFields(o.payload.fields)
}

// Names returns detached field names in original order.
//
// Names mirrors Fields order and returns a caller-owned slice.
func (o ObjectView) Names() []string {
	names := make([]string, 0, len(o.payload.fields))
	for _, field := range o.payload.fields {
		names = append(names, field.Name)
	}

	return names
}

// findObjectField returns the ordered field index for name.
//
// A negative result means the field is absent. Keeping the search as a tiny
// helper makes Has and Get share the same linear lookup semantics.
func findObjectField(fields []Field, name string) int {
	for i, field := range fields {
		if field.Name == name {
			return i
		}
	}

	return -1
}
