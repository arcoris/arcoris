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

// ListValue constructs a KindList Value from initialized item values.
//
// Lists preserve caller order and do not store element descriptors, uniqueness
// policy, or length constraints. Each item is cloned so the caller keeps no
// mutable alias into the stored list.
func ListValue(items ...Value) (Value, error) {
	payload, err := newListPayload(items)
	if err != nil {
		return Value{}, err
	}

	return Value{kind: KindList, listValue: payload}, nil
}

// MustListValue constructs a list Value or panics when items are malformed.
//
// It is intended for tests and static fixtures. Runtime construction paths
// should use ListValue and handle its structured error.
func MustListValue(items ...Value) Value {
	value, err := ListValue(items...)
	if err != nil {
		panic(err)
	}

	return value
}
