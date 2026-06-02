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

package listmapkey

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/value"
)

// build extracts and validates the selector for the configured item.
func (r selectorRequest) build() (fieldpath.Selector, error) {
	if len(r.keys) == 0 {
		return fieldpath.Selector{}, failure(
			r.indexPath,
			FailureInvalidDescriptor,
			"ListMap has no key fields",
		)
	}

	if r.item.Kind() != value.KindObject {
		return fieldpath.Selector{}, failure(
			r.indexPath,
			FailureItemKindMismatch,
			"ListMap item is not an object",
		)
	}

	elementObjectDescriptor, err := r.objectDescriptor()
	if err != nil {
		return fieldpath.Selector{}, err
	}

	itemObjectView, _ := r.item.Object()
	selectorEntries := make([]fieldpath.SelectorEntry, 0, len(r.keys))

	for _, key := range r.keys {
		selectorEntryValue, err := r.selectorEntry(
			itemObjectView,
			elementObjectDescriptor,
			key,
		)
		if err != nil {
			return fieldpath.Selector{}, err
		}

		selectorEntries = append(selectorEntries, selectorEntryValue)
	}

	itemSelector, err := fieldpath.NewSelector(selectorEntries...)
	if err != nil {
		return fieldpath.Selector{}, failureWithCause(
			r.indexPath,
			FailureInvalidSelector,
			"ListMap selector is invalid",
			err,
		)
	}

	return itemSelector, nil
}
