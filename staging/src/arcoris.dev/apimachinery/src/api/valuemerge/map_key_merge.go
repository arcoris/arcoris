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

// mergeBaseMapMembers walks base keys, preserving unselected keys.
func (m *merger) mergeBaseMapMembers(
	path fieldpath.Path,
	baseMembers []value.RecordMember,
	overlayLookup memberLookup,
	valueDescriptor types.Descriptor,
	fields fieldpath.Set,
	depth int,
) ([]value.RecordMember, error) {
	members := make([]value.RecordMember, 0, len(baseMembers))

	for _, member := range baseMembers {
		name := member.Name.String()
		childPath := path.Key(fieldpath.MustMapKey(name))
		next, err := m.mergeMapMember(
			childPath,
			valuepresence.Present(member.Value),
			overlayLookup.Operand(name),
			valueDescriptor,
			fields,
			depth,
		)
		if err != nil {
			return nil, err
		}

		members = appendMember(members, name, next)
	}

	return members, nil
}

// appendOverlayMapMembers appends selected overlay keys absent from base.
func (m *merger) appendOverlayMapMembers(
	path fieldpath.Path,
	members []value.RecordMember,
	baseLookup memberLookup,
	overlayMembers []value.RecordMember,
	valueDescriptor types.Descriptor,
	fields fieldpath.Set,
	depth int,
) ([]value.RecordMember, error) {
	for _, member := range overlayMembers {
		name := member.Name.String()
		if baseLookup.Has(name) {
			continue
		}

		childPath := path.Key(fieldpath.MustMapKey(name))
		next, err := m.mergeMapMember(
			childPath,
			valuepresence.Absent(),
			valuepresence.Present(member.Value),
			valueDescriptor,
			fields,
			depth,
		)
		if err != nil {
			return nil, err
		}

		members = appendMember(members, name, next)
	}

	return members, nil
}

// mergeMapMember merges one dynamic key when selected.
func (m *merger) mergeMapMember(
	path fieldpath.Path,
	base operand,
	overlay operand,
	valueDescriptor types.Descriptor,
	fields fieldpath.Set,
	depth int,
) (operand, error) {
	if !hasSelectedChild(fields, path) {
		return base.Clone(), nil
	}

	return m.merge(path, base, overlay, valueDescriptor, fields, depth+1)
}
