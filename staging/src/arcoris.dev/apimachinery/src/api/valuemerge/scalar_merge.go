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
	"arcoris.dev/apimachinery/api/internal/typekind"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// mergeScalar rejects selected descendants below scalar descriptors.
func (m *merger) mergeScalar(
	path fieldpath.Path,
	base operand,
	overlay operand,
	descriptor types.Type,
) (operand, error) {
	if err := requireScalarOperands(path, base, overlay, descriptor); err != nil {
		return operand{}, err
	}

	return operand{}, errorfAt(
		path,
		ErrUnsupportedMerge,
		ErrorReasonUnsupportedMerge,
		"descriptor %s has no mergeable descendants",
		descriptorKindName(descriptor),
	)
}

// requireScalarOperands validates the visible scalar operands before reporting
// descendant-selection failure. Exact scalar replacement is handled by merge.
func requireScalarOperands(
	path fieldpath.Path,
	base operand,
	overlay operand,
	descriptor types.Type,
) error {
	expected, ok := scalarKind(descriptor.Code())
	if !ok {
		return errorfAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor %s is not mergeable here",
			descriptor.Code(),
		)
	}
	if err := requireKind(path, base, expected); err != nil {
		return err
	}

	return requireKind(path, overlay, expected)
}

// scalarKind maps scalar descriptor codes to concrete payload kinds.
func scalarKind(code types.TypeCode) (value.Kind, bool) {
	if code == types.TypeNull {
		return value.KindNull, true
	}

	return typekind.Scalar(code)
}
