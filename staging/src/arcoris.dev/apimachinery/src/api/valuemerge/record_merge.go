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

// mergeRecord recursively merges record payload members under a fixed object descriptor.
func (m *merger) mergeRecord(
	path fieldpath.Path,
	base operand,
	overlay operand,
	descriptor types.Descriptor,
	fields fieldpath.Set,
	depth int,
) (operand, error) {
	if err := requireObjectOperand(path, base); err != nil {
		return operand{}, err
	}
	if err := requireObjectOperand(path, overlay); err != nil {
		return operand{}, err
	}
	preserved, ok, err := preserveWithoutOverlayContainer(path, base, overlay, value.KindRecord)
	if err != nil {
		return operand{}, err
	}
	if ok {
		return preserved, nil
	}

	view, ok := descriptor.AsObject()
	if !ok {
		return operand{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not an object",
		)
	}

	return m.mergeRecordView(path, base, overlay, view, fields, depth)
}

// mergeRecordView merges one record payload using its object descriptor view.
func (m *merger) mergeRecordView(
	path fieldpath.Path,
	base operand,
	overlay operand,
	view types.ObjectView,
	fields fieldpath.Set,
	depth int,
) (operand, error) {
	baseMembers := recordMembers(base)
	overlayMembers := recordMembers(overlay)
	baseLookup := newMemberLookup(baseMembers)
	overlayLookup := newMemberLookup(overlayMembers)
	declared := newRecordFieldLookup(view.Fields())

	members, err := m.mergeBaseRecordMembers(
		path,
		baseMembers,
		baseLookup,
		overlayLookup,
		declared,
		view.UnknownFields(),
		fields,
		depth,
	)
	if err != nil {
		return operand{}, err
	}

	members, err = m.appendOverlayRecordMembers(
		path,
		members,
		baseLookup,
		overlayMembers,
		declared,
		view.UnknownFields(),
		fields,
		depth,
	)
	if err != nil {
		return operand{}, err
	}

	merged, err := value.RecordValue(members...)
	if err != nil {
		return operand{}, wrapAt(
			path,
			ErrInvalidValue,
			ErrorReasonInvalidMergedValue,
			"merged record is invalid",
			err,
		)
	}

	return valuepresence.Present(merged), nil
}
