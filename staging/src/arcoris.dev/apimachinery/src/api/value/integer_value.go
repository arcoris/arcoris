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

// Int64 constructs an integer Value from a signed 64-bit value.
//
// The resulting payload kind is always KindInteger; signedness is payload data,
// not a separate kind.
func Int64(v int64) Value {
	return IntegerValue(NewIntegerFromInt64(v))
}

// Uint64 constructs an integer Value from an unsigned 64-bit value.
//
// The resulting payload kind is always KindInteger. Width-specific constraints
// belong to descriptors and future value validation.
func Uint64(v uint64) Value {
	return IntegerValue(NewIntegerFromUint64(v))
}

// IntegerValue constructs an integer Value from an already canonical Integer.
//
// The Integer type is immutable by convention, so the value can be copied
// directly into the active payload slot.
func IntegerValue(v Integer) Value {
	if v.magnitude == 0 {
		v.negative = false
	}

	return Value{kind: KindInteger, integerValue: v}
}
