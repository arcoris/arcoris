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
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// equalValue reports whether two present concrete values are equal under
// descriptor semantics.
//
// It is used only for leaf-style decisions: scalar modification checks,
// atomic/list-set whole-field checks, and descriptor-aware structural equality.
// Added and removed subtrees still go through valuefieldset so their paths match
// ownership/field-set semantics.
func (c *comparer) equalValue(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
	depth int,
) (bool, error) {
	if oldValue.IsZero() || newValue.IsZero() {
		return false, errorAt(
			path,
			ErrInvalidValue,
			ErrorReasonInvalidZero,
			"value is the invalid zero Value",
		)
	}
	if !descriptor.IsValid() {
		return false, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor has no valid type code",
		)
	}

	if oldValue.IsNull() || newValue.IsNull() {
		return oldValue.IsNull() && newValue.IsNull(), nil
	}

	switch descriptor.Code() {
	case types.TypeNull:
		if err := requireKind(path, oldValue, value.KindNull, descriptor.Code()); err != nil {
			return false, err
		}
		if err := requireKind(path, newValue, value.KindNull, descriptor.Code()); err != nil {
			return false, err
		}
		return true, nil
	case types.TypeBool,
		types.TypeString,
		types.TypeBytes,
		types.TypeInt8,
		types.TypeInt16,
		types.TypeInt32,
		types.TypeInt64,
		types.TypeUint8,
		types.TypeUint16,
		types.TypeUint32,
		types.TypeUint64,
		types.TypeFloat32,
		types.TypeFloat64,
		types.TypeDecimal,
		types.TypeTimestamp,
		types.TypeDate,
		types.TypeTime,
		types.TypeDuration:
		return c.equalScalar(path, oldValue, newValue, descriptor)
	case types.TypeObject:
		return c.equalObject(path, oldValue, newValue, descriptor, depth)
	case types.TypeMap:
		return c.equalMap(path, oldValue, newValue, descriptor, depth)
	case types.TypeList:
		return c.equalList(path, oldValue, newValue, descriptor, depth)
	case types.TypeRef:
		name, resolved, err := c.resolveRefDefinition(path, descriptor, depth)
		if err != nil {
			return false, err
		}

		c.resolving[name] = true
		defer delete(c.resolving, name)

		return c.equalValue(path, oldValue, newValue, resolved, depth+1)
	default:
		return false, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor has an unsupported type code",
		)
	}
}
