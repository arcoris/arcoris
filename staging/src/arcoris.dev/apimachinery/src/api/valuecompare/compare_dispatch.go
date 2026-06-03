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
)

// compare applies presence rules first, then dispatches by descriptor code.
func (c *comparer) compare(
	path fieldpath.Path,
	oldOperand operand,
	newOperand operand,
	descriptor types.Type,
	depth int,
) (Result, error) {
	if result, done, err := c.comparePresence(path, oldOperand, newOperand, descriptor); done {
		return result, err
	}

	oldValue := oldOperand.value
	newValue := newOperand.value
	if err := requireComparableInputs(path, oldValue, newValue, descriptor); err != nil {
		return Result{}, err
	}
	if oldValue.IsNull() || newValue.IsNull() {
		return c.compareNull(path, oldValue, newValue)
	}

	switch descriptor.Code() {
	case types.TypeNull:
		return c.compareNullDescriptor(path, oldValue, newValue, descriptor)
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
		return c.compareScalar(path, oldValue, newValue, descriptor)
	case types.TypeObject:
		return c.compareObject(path, oldValue, newValue, descriptor, depth)
	case types.TypeMap:
		return c.compareMap(path, oldValue, newValue, descriptor, depth)
	case types.TypeList:
		return c.compareList(path, oldValue, newValue, descriptor, depth)
	case types.TypeRef:
		return c.compareRef(path, oldOperand, newOperand, descriptor, depth)
	default:
		return Result{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor has an unsupported type code",
		)
	}
}
