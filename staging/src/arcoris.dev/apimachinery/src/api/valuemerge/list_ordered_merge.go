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

// mergeOrderedList merges physical indexes for ordered-list descriptors.
func (m *merger) mergeOrderedList(
	path fieldpath.Path,
	base operand,
	overlay operand,
	element types.Type,
	fields fieldpath.Set,
	depth int,
) (operand, error) {
	baseItems := listItems(base)
	overlayItems := listItems(overlay)
	items := make([]value.Value, 0, orderedListCapacity(baseItems, overlayItems))

	for i, item := range baseItems {
		childPath := path.Index(i)
		next, err := m.mergeOrderedItem(
			childPath,
			valuepresence.Present(item),
			itemAt(overlayItems, i),
			element,
			fields,
			depth,
		)
		if err != nil {
			return operand{}, err
		}

		items = appendItem(items, next)
	}

	for i := len(baseItems); i < len(overlayItems); i++ {
		childPath := path.Index(i)
		next, err := m.mergeOrderedItem(
			childPath,
			valuepresence.Absent(),
			valuepresence.Present(overlayItems[i]),
			element,
			fields,
			depth,
		)
		if err != nil {
			return operand{}, err
		}

		items = appendItem(items, next)
	}

	merged, err := value.ListValue(items...)
	if err != nil {
		return operand{}, wrapAt(
			path,
			ErrInvalidValue,
			ErrorReasonInvalidZero,
			"merged list is invalid",
			err,
		)
	}

	return valuepresence.Present(merged), nil
}

// mergeOrderedItem merges one selected ordered-list index.
func (m *merger) mergeOrderedItem(
	path fieldpath.Path,
	base operand,
	overlay operand,
	element types.Type,
	fields fieldpath.Set,
	depth int,
) (operand, error) {
	if !hasSelectedChild(fields, path) {
		return base.Clone(), nil
	}

	return m.merge(path, base, overlay, element, fields, depth+1)
}

// orderedListCapacity picks enough space for preserved and appended items.
func orderedListCapacity(base []value.Value, overlay []value.Value) int {
	if len(overlay) > len(base) {
		return len(overlay)
	}

	return len(base)
}
