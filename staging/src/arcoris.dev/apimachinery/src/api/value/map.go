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

// NewMap constructs a map value from initialized, uniquely keyed entries.
//
// Map values represent dynamic key/value payloads and are distinct from object
// records. They do not store key descriptors, value descriptors, key regexes,
// or entry-count constraints. Entries are cloned and kept in caller order.
func NewMap(entries ...Entry) (Value, error) {
	payload, err := newMapPayload(entries)
	if err != nil {
		return Value{}, err
	}

	return Value{kind: KindMap, mapValue: payload}, nil
}

// MustMap constructs a map Value or panics when entries are malformed.
//
// It is intended for tests and static fixtures. Runtime construction paths
// should use NewMap and handle its structured error.
func MustMap(entries ...Entry) Value {
	value, err := NewMap(entries...)
	if err != nil {
		panic(err)
	}

	return value
}
