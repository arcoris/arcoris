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

// equalValue reports whether two present payload values are equal by descriptor semantics.
//
// It is used for whole-value decisions that need equality but not a diff set:
// scalar changes, atomic/list-set list comparison, opaque-preserve comparison,
// and nested equality under those cases. Added and removed subtrees still go
// through valuefieldset so their paths match field-set semantics.
func (c *comparer) equalValue(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Descriptor,
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
			"descriptor has no valid kind",
		)
	}

	if oldValue.IsNull() || newValue.IsNull() {
		return oldValue.IsNull() && newValue.IsNull(), nil
	}

	switch descriptor.Code() {
	case types.DescriptorNull:
		if err := requireKind(path, oldValue, value.KindNull, descriptor.Code()); err != nil {
			return false, err
		}
		if err := requireKind(path, newValue, value.KindNull, descriptor.Code()); err != nil {
			return false, err
		}
		return true, nil
	case types.DescriptorBool,
		types.DescriptorString,
		types.DescriptorBytes,
		types.DescriptorInt8,
		types.DescriptorInt16,
		types.DescriptorInt32,
		types.DescriptorInt64,
		types.DescriptorUint8,
		types.DescriptorUint16,
		types.DescriptorUint32,
		types.DescriptorUint64,
		types.DescriptorFloat32,
		types.DescriptorFloat64,
		types.DescriptorDecimal,
		types.DescriptorTimestamp,
		types.DescriptorDate,
		types.DescriptorTime,
		types.DescriptorDuration:
		return c.equalScalar(path, oldValue, newValue, descriptor)
	case types.DescriptorObject:
		return c.equalObject(path, oldValue, newValue, descriptor, depth)
	case types.DescriptorMap:
		return c.equalMap(path, oldValue, newValue, descriptor, depth)
	case types.DescriptorList:
		return c.equalList(path, oldValue, newValue, descriptor, depth)
	case types.DescriptorRef:
		name, resolved, err := c.resolveRefDefinition(path, descriptor, depth)
		if err != nil {
			return false, err
		}

		leave := c.refs.Enter(name)
		defer leave()

		return c.equalValue(path, oldValue, newValue, resolved, depth+1)
	default:
		return false, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor has an unsupported kind",
		)
	}
}
