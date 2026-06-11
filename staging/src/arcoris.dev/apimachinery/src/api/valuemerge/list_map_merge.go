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

// mergeListMap merges associative-list items by selector identity.
func (m *merger) mergeListMap(
	path fieldpath.Path,
	base operand,
	overlay operand,
	element types.Descriptor,
	keys []types.FieldName,
	fields fieldpath.Set,
	depth int,
) (operand, error) {
	baseIndex, err := m.listMapIndex(path, base, element, keys)
	if err != nil {
		return operand{}, err
	}

	overlayIndex, err := m.listMapIndex(path, overlay, element, keys)
	if err != nil {
		return operand{}, err
	}

	items, err := m.mergeBaseListMapItems(path, baseIndex, overlayIndex, element, fields, depth)
	if err != nil {
		return operand{}, err
	}

	items, err = m.appendOverlayListMapItems(path, items, baseIndex, overlayIndex, element, fields, depth)
	if err != nil {
		return operand{}, err
	}

	merged, err := value.ListValue(items...)
	if err != nil {
		return operand{}, wrapAt(
			path,
			ErrInvalidValue,
			ErrorReasonInvalidMergedValue,
			"merged list map is invalid",
			err,
		)
	}

	return valuepresence.Present(merged), nil
}

// mergeBaseListMapItems preserves base order for existing selectors.
func (m *merger) mergeBaseListMapItems(
	path fieldpath.Path,
	baseIndex listMapIndex,
	overlayIndex listMapIndex,
	element types.Descriptor,
	fields fieldpath.Set,
	depth int,
) ([]value.Value, error) {
	items := make([]value.Value, 0, len(baseIndex.order)+len(overlayIndex.order))

	for _, key := range baseIndex.order {
		baseEntry := baseIndex.lookup[key]
		overlayEntry, overlayFound := overlayIndex.lookup[key]
		selectorPath := path.Select(baseEntry.selector)

		next, err := m.mergeListMapItem(
			selectorPath,
			valuepresence.Present(baseEntry.item),
			listMapOperand(overlayEntry, overlayFound),
			element,
			fields,
			depth,
		)
		if err != nil {
			return nil, err
		}

		items = appendItem(items, next)
	}

	return items, nil
}

// appendOverlayListMapItems appends selected overlay selectors absent from base.
func (m *merger) appendOverlayListMapItems(
	path fieldpath.Path,
	items []value.Value,
	baseIndex listMapIndex,
	overlayIndex listMapIndex,
	element types.Descriptor,
	fields fieldpath.Set,
	depth int,
) ([]value.Value, error) {
	for _, key := range overlayIndex.order {
		if _, exists := baseIndex.lookup[key]; exists {
			continue
		}

		entry := overlayIndex.lookup[key]
		selectorPath := path.Select(entry.selector)
		next, err := m.mergeListMapItem(
			selectorPath,
			valuepresence.Absent(),
			valuepresence.Present(entry.item),
			element,
			fields,
			depth,
		)
		if err != nil {
			return nil, err
		}

		items = appendItem(items, next)
	}

	return items, nil
}

// mergeListMapItem merges one selector-addressed item when selected.
func (m *merger) mergeListMapItem(
	path fieldpath.Path,
	base operand,
	overlay operand,
	element types.Descriptor,
	fields fieldpath.Set,
	depth int,
) (operand, error) {
	if !hasSelectedChild(fields, path) {
		return base.Clone(), nil
	}

	return m.merge(path, base, overlay, element, fields, depth+1)
}

// listMapOperand converts a selector lookup into presence-aware merge input.
func listMapOperand(entry listMapEntry, present bool) operand {
	return valuepresence.From(entry.item, present)
}
