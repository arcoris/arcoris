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
	"arcoris.dev/apimachinery/api/value"
)

// replaceSubtree copies overlay after the minimal shape checks needed for merge.
func (m *merger) replaceSubtree(
	path fieldpath.Path,
	overlay operand,
	descriptor types.Descriptor,
	depth int,
) (operand, error) {
	if overlay.Absent() {
		return valuepresence.Absent(), nil
	}
	if err := m.requireReplacementKind(path, overlay, descriptor, depth); err != nil {
		return operand{}, err
	}

	return overlay.Clone(), nil
}

// requireReplacementKind checks only zero value, DescriptorRef, and descriptor kind.
func (m *merger) requireReplacementKind(
	path fieldpath.Path,
	overlay operand,
	descriptor types.Descriptor,
	depth int,
) error {
	if err := requireValidValue(path, overlay); err != nil {
		return err
	}
	if overlay.Absent() || overlay.Value().IsNull() {
		return nil
	}

	switch descriptor.Code() {
	case types.DescriptorRef:
		name, resolved, err := m.resolveRefDefinition(path, descriptor, depth)
		if err != nil {
			return err
		}

		leave := m.refs.Enter(name)
		defer leave()

		return m.requireReplacementKind(path, overlay, resolved, depth+1)
	case types.DescriptorObject, types.DescriptorMap:
		return requireKind(path, overlay, value.KindObject)
	case types.DescriptorList:
		return requireKind(path, overlay, value.KindList)
	default:
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

		return requireKind(path, overlay, expected)
	}
}
