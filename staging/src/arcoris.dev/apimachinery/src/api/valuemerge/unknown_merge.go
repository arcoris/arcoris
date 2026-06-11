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
	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/types"
)

// mergeUnknownRecordMember handles undeclared record members as policy leaves.
func mergeUnknownRecordMember(
	path fieldpath.Path,
	base operand,
	overlay operand,
	unknown types.UnknownFieldPolicy,
	selection mergeSelection,
) (operand, error) {
	switch unknown {
	case types.UnknownReject:
		return operand{}, errorAt(
			path,
			ErrUnknownField,
			ErrorReasonUnknownField,
			"record member is not declared by the object descriptor",
		)
	case types.UnknownPrune:
		return valuepresence.Absent(), nil
	case types.UnknownPreserveOpaque:
		return mergePreservedUnknownMember(path, base, overlay, selection)
	default:
		return operand{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"unknown-field policy is invalid",
		)
	}
}

// mergePreservedUnknownMember treats an undeclared field as one opaque leaf.
func mergePreservedUnknownMember(
	path fieldpath.Path,
	base operand,
	overlay operand,
	selection mergeSelection,
) (operand, error) {
	if !selection.exact && !selection.descendants.IsEmpty() {
		return operand{}, errorAt(
			path,
			ErrUnsupportedMerge,
			ErrorReasonUnsupportedMerge,
			"cannot merge below preserved unknown field",
		)
	}
	if !selection.exact {
		return base.Clone(), nil
	}
	if overlay.Absent() {
		return valuepresence.Absent(), nil
	}

	return overlay.Clone(), nil
}
