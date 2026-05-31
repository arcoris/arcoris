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

// NewObject constructs an object value from initialized, uniquely named fields.
//
// Object values represent schema-shaped records with caller-ordered fields.
// They do not store descriptor metadata such as required fields, optional
// fields, JSON tags, or unknown-field policy. Field values are cloned during
// construction so later mutations of caller-owned values cannot affect the
// stored object.
func NewObject(fields ...Field) (Value, error) {
	payload, err := newObjectPayload(fields)
	if err != nil {
		return Value{}, err
	}

	return Value{kind: KindObject, objectValue: payload}, nil
}

// MustObject constructs an object Value or panics when fields are malformed.
//
// It is intended for tests and static fixtures where malformed field data is a
// programmer error. Runtime construction paths should use NewObject and return
// its structured error.
func MustObject(fields ...Field) Value {
	value, err := NewObject(fields...)
	if err != nil {
		panic(err)
	}

	return value
}
