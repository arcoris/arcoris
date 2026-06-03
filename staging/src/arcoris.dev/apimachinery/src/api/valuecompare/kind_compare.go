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

// requireComparableInputs rejects invalid zero values and invalid descriptors.
func requireComparableInputs(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
) error {
	if oldValue.IsZero() || newValue.IsZero() {
		return errorAt(
			path,
			ErrInvalidValue,
			ErrorReasonInvalidZero,
			"value is the invalid zero Value",
		)
	}
	if !descriptor.IsValid() {
		return errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor has no valid type code",
		)
	}

	return nil
}

// requireKind reports a concrete kind / descriptor type mismatch.
func requireKind(path fieldpath.Path, val value.Value, expected value.Kind, code types.TypeCode) error {
	if val.Kind() == expected {
		return nil
	}

	return errorfAt(
		path,
		ErrKindMismatch,
		ErrorReasonKindMismatch,
		"value kind %s does not match descriptor %s; expected %s",
		val.Kind(),
		code,
		expected,
	)
}
