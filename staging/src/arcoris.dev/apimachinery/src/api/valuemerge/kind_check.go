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

package valuemerge

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// requireValidValue rejects the invalid zero value when an operand is present.
func requireValidValue(path fieldpath.Path, operand operand) error {
	if operand.Absent() || !operand.Value().IsZero() {
		return nil
	}

	return errorAt(
		path,
		ErrInvalidValue,
		ErrorReasonInvalidZero,
		"value is invalid",
	)
}

// requireKind checks a present non-null value against the expected concrete kind.
func requireKind(path fieldpath.Path, operand operand, expected value.Kind) error {
	if err := requireValidValue(path, operand); err != nil {
		return err
	}
	if operand.Absent() || operand.Value().IsNull() {
		return nil
	}
	if operand.Value().Kind() == expected {
		return nil
	}

	return errorfAt(
		path,
		ErrKindMismatch,
		ErrorReasonKindMismatch,
		"value kind %s does not match descriptor kind %s",
		operand.Value().Kind(),
		expected,
	)
}

// requireObjectOperand checks object/map concrete payload shape.
func requireObjectOperand(path fieldpath.Path, operand operand) error {
	return requireKind(path, operand, value.KindObject)
}

// requireListOperand checks list concrete payload shape.
func requireListOperand(path fieldpath.Path, operand operand) error {
	return requireKind(path, operand, value.KindList)
}

// descriptorKindName returns a stable diagnostic descriptor name.
func descriptorKindName(descriptor types.Descriptor) string {
	if descriptor.Code() == types.DescriptorNull {
		return value.KindNull.String()
	}

	return descriptor.Code().String()
}
