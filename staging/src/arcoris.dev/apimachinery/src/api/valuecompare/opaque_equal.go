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

package valuecompare

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/value"
)

// equalOpaqueValue compares preserved unknown fields without descriptor descent.
//
// Unknown-preserved members are semantic leaves for comparison output. Their
// nested payload is inspected only to decide whether the leaf changed.
func (c *comparer) equalOpaqueValue(path fieldpath.Path, oldValue value.Value, newValue value.Value) (bool, error) {
	if oldValue.IsZero() || newValue.IsZero() {
		return false, errorAt(path, ErrInvalidValue, ErrorReasonInvalidZero, "value is the invalid zero Value")
	}
	if oldValue.Kind() != newValue.Kind() {
		return false, nil
	}

	if oldValue.Kind().IsPrimitive() {
		return opaqueScalarValuesEqual(oldValue, newValue), nil
	}

	switch oldValue.Kind() {
	case value.KindObject:
		return c.equalOpaqueObject(path, oldValue, newValue)
	case value.KindList:
		return c.equalOpaqueList(path, oldValue, newValue)
	default:
		return false, errorAt(path, ErrInvalidValue, ErrorReasonInvalidZero, "value has invalid kind")
	}
}
