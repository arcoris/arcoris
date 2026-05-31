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

// ObjectValue constructs a KindObject Value from uniquely named members.
//
// Object values are concrete keyed payload nodes with caller-ordered members.
// They intentionally do not decide whether the data will satisfy a descriptor
// object or descriptor map. Member values are cloned during construction so
// later mutations of caller-owned values cannot affect the stored object.
func ObjectValue(members ...Member) (Value, error) {
	payload, err := newObjectPayload(members)
	if err != nil {
		return Value{}, err
	}

	return Value{kind: KindObject, objectValue: payload}, nil
}

// MustObjectValue constructs an object Value or panics when members are malformed.
//
// It is intended for tests and static fixtures where malformed member data is a
// programmer error. Runtime construction paths should use ObjectValue and return
// its structured error.
func MustObjectValue(members ...Member) Value {
	value, err := ObjectValue(members...)
	if err != nil {
		panic(err)
	}

	return value
}
