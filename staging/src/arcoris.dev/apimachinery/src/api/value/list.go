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

// NewList constructs a list value from initialized item values.
//
// Lists preserve caller order and do not store element descriptors, uniqueness
// policy, set/list/map semantics, or length constraints. Each item is cloned so
// the caller keeps no mutable alias into the stored list.
func NewList(items ...Value) (Value, error) {
	payload, err := newListPayload(items)
	if err != nil {
		return Value{}, err
	}

	return Value{kind: KindList, listValue: payload}, nil
}

// MustList constructs a list Value or panics when items are malformed.
//
// It is intended for tests and static fixtures. Runtime construction paths
// should use NewList and handle its structured error.
func MustList(items ...Value) Value {
	value, err := NewList(items...)
	if err != nil {
		panic(err)
	}

	return value
}
