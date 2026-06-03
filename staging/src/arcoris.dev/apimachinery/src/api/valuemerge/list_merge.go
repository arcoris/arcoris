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

// mergeList dispatches by descriptor list semantics.
func (m *merger) mergeList(
	path fieldpath.Path,
	base operand,
	overlay operand,
	descriptor types.Type,
	fields fieldpath.Set,
	depth int,
) (operand, error) {
	if err := requireListOperand(path, base); err != nil {
		return operand{}, err
	}
	if err := requireListOperand(path, overlay); err != nil {
		return operand{}, err
	}
	if preserved, ok := preserveWithoutOverlayContainer(base, overlay, value.KindList); ok {
		return preserved, nil
	}

	view, ok := descriptor.List()
	if !ok {
		return operand{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not a list",
		)
	}

	switch view.Semantics() {
	case types.ListAtomic, types.ListSet:
		return operand{}, errorAt(
			path,
			ErrUnsupportedMerge,
			ErrorReasonUnsupportedMerge,
			"cannot merge below atomic or set list",
		)
	case types.ListOrdered:
		return m.mergeOrderedList(path, base, overlay, view.Element(), fields, depth)
	case types.ListMap:
		return m.mergeListMap(path, base, overlay, view.Element(), view.MapKeys(), fields, depth)
	default:
		return operand{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"list semantics are invalid",
		)
	}
}
