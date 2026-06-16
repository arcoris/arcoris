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

package objectquery

import "arcoris.dev/apimachinery/api/value"

// compareOrdered compares same-kind values that the compiler has classified as
// orderable. The bool result is false for unsupported runtime kinds.
func compareOrdered(left value.Value, right value.Value) (int, bool) {
	if left.Kind() != right.Kind() {
		return 0, false
	}

	switch left.Kind() {
	case value.KindInteger:
		l, _ := left.AsInteger()
		r, _ := right.AsInteger()
		return l.Compare(r), true
	case value.KindFloat:
		return compareFloats(left, right)
	case value.KindDecimal:
		l, _ := left.AsDecimal()
		r, _ := right.AsDecimal()
		return l.Compare(r), true
	case value.KindTimestamp:
		l, _ := left.AsTimestamp()
		r, _ := right.AsTimestamp()
		return l.Compare(r), true
	case value.KindDate, value.KindTimeOfDay, value.KindDuration:
		return compareCanonicalKeys(left, right), true
	default:
		return 0, false
	}
}

// compareFloats returns a three-way comparison for float values.
func compareFloats(left value.Value, right value.Value) (int, bool) {
	l, _ := left.AsFloat()
	r, _ := right.AsFloat()
	switch {
	case l < r:
		return -1, true
	case l > r:
		return 1, true
	default:
		return 0, true
	}
}

// compareCanonicalKeys provides stable ordering for value kinds whose canonical
// string form is already their ordered representation in this package.
func compareCanonicalKeys(left value.Value, right value.Value) int {
	l := canonicalValueKey(left)
	r := canonicalValueKey(right)
	switch {
	case l < r:
		return -1
	case l > r:
		return 1
	default:
		return 0
	}
}
