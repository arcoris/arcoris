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

// mergeBaseObjectMembers walks base order, preserving unselected fields.
func (m *merger) mergeBaseObjectMembers(
	path fieldpath.Path,
	baseMembers []value.Member,
	baseLookup memberLookup,
	overlayLookup memberLookup,
	declared objectFieldLookup,
	unknown types.UnknownFieldPolicy,
	fields fieldpath.Set,
	depth int,
) ([]value.Member, error) {
	members := make([]value.Member, 0, len(baseMembers))

	for _, member := range baseMembers {
		field, known := declared[member.Name]
		childPath := path.Field(member.Name)

		next, err := m.mergeObjectMember(
			childPath,
			baseLookup.Operand(member.Name),
			overlayLookup.Operand(member.Name),
			field,
			known,
			unknown,
			fields,
			depth,
		)
		if err != nil {
			return nil, err
		}

		members = appendMember(members, member.Name, next)
	}

	return members, nil
}

// appendOverlayObjectMembers appends selected overlay members absent from base.
func (m *merger) appendOverlayObjectMembers(
	path fieldpath.Path,
	members []value.Member,
	baseLookup memberLookup,
	overlayMembers []value.Member,
	declared objectFieldLookup,
	unknown types.UnknownFieldPolicy,
	fields fieldpath.Set,
	depth int,
) ([]value.Member, error) {
	for _, member := range overlayMembers {
		if baseLookup.Has(member.Name) {
			continue
		}

		field, known := declared[member.Name]
		childPath := path.Field(member.Name)
		next, err := m.mergeObjectMember(
			childPath,
			valuepresence.Absent(),
			valuepresence.Present(member.Value),
			field,
			known,
			unknown,
			fields,
			depth,
		)
		if err != nil {
			return nil, err
		}

		members = appendMember(members, member.Name, next)
	}

	return members, nil
}

// mergeObjectMember applies object unknown-field policy before recursion.
func (m *merger) mergeObjectMember(
	path fieldpath.Path,
	base operand,
	overlay operand,
	field types.FieldDescriptor,
	known bool,
	unknown types.UnknownFieldPolicy,
	fields fieldpath.Set,
	depth int,
) (operand, error) {
	selection := selectAt(fields, path)
	if !selection.selected() {
		switch {
		case known || unknown == types.UnknownPreserve:
			return base.Clone(), nil
		case unknown == types.UnknownReject:
			return operand{}, errorAt(
				path,
				ErrUnknownField,
				ErrorReasonUnknownField,
				"object field is not declared",
			)
		default:
			return valuepresence.Absent(), nil
		}
	}

	if known {
		return m.merge(path, base, overlay, field.Type(), fields, depth+1)
	}

	return mergeUnknownObjectMember(path, base, overlay, unknown, selection)
}
