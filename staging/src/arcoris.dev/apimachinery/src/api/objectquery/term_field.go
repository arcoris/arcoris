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

// FieldExists matches items where ref resolves to a present value.
func FieldExists(ref FieldRef) (Query, error) {
	return fieldQuery(ref, OperatorExists)
}

// FieldDoesNotExist matches items where ref is absent.
func FieldDoesNotExist(ref FieldRef) (Query, error) {
	return fieldQuery(ref, OperatorDoesNotExist)
}

// FieldEquals matches items where ref is present and equals literal.
func FieldEquals(ref FieldRef, literal value.Value) (Query, error) {
	return fieldQuery(ref, OperatorEquals, literal)
}

// FieldNotEquals matches items where ref is absent or differs from literal.
func FieldNotEquals(ref FieldRef, literal value.Value) (Query, error) {
	return fieldQuery(ref, OperatorNotEquals, literal)
}

// FieldIn matches items where ref is present and belongs to literals.
func FieldIn(ref FieldRef, literals ...value.Value) (Query, error) {
	return fieldQuery(ref, OperatorIn, literals...)
}

// FieldNotIn matches items where ref is absent or outside literals.
func FieldNotIn(ref FieldRef, literals ...value.Value) (Query, error) {
	return fieldQuery(ref, OperatorNotIn, literals...)
}

// FieldLessThan matches items where ref is less than literal.
func FieldLessThan(ref FieldRef, literal value.Value) (Query, error) {
	return fieldQuery(ref, OperatorLessThan, literal)
}

// FieldLessOrEqual matches items where ref is less than or equal to literal.
func FieldLessOrEqual(ref FieldRef, literal value.Value) (Query, error) {
	return fieldQuery(ref, OperatorLessOrEqual, literal)
}

// FieldGreaterThan matches items where ref is greater than literal.
func FieldGreaterThan(ref FieldRef, literal value.Value) (Query, error) {
	return fieldQuery(ref, OperatorGreaterThan, literal)
}

// FieldGreaterOrEqual matches items where ref is greater than or equal to literal.
func FieldGreaterOrEqual(ref FieldRef, literal value.Value) (Query, error) {
	return fieldQuery(ref, OperatorGreaterOrEqual, literal)
}

// FieldHasPrefix matches string fields with prefix.
func FieldHasPrefix(ref FieldRef, prefix string) (Query, error) {
	return fieldQuery(ref, OperatorHasPrefix, value.StringValue(prefix))
}

// FieldHasSuffix matches string fields with suffix.
func FieldHasSuffix(ref FieldRef, suffix string) (Query, error) {
	return fieldQuery(ref, OperatorHasSuffix, value.StringValue(suffix))
}

// FieldContains matches string fields containing needle.
func FieldContains(ref FieldRef, needle string) (Query, error) {
	return fieldQuery(ref, OperatorContains, value.StringValue(needle))
}
