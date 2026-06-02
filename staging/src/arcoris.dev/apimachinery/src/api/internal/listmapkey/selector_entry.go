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
	"fmt"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// selectorEntry extracts one declared ListMap key field into a selector entry.
func (r selectorRequest) selectorEntry(
	itemObjectView value.ObjectView,
	elementObjectDescriptor types.ObjectView,
	key types.FieldName,
) (fieldpath.SelectorEntry, error) {
	keyName := string(key)
	keyPath := r.indexPath.Field(keyName)

	keyFieldDescriptor, ok := fieldDescriptor(elementObjectDescriptor, key)
	if !ok {
		return fieldpath.SelectorEntry{}, failure(
			keyPath,
			FailureInvalidDescriptor,
			fmt.Sprintf(
				"ListMap key field %q is not declared by the element descriptor",
				keyName,
			),
		)
	}

	keyMemberValue, ok := itemObjectView.Get(keyName)
	if !ok {
		return fieldpath.SelectorEntry{}, failure(
			keyPath,
			FailureMissingKey,
			fmt.Sprintf("ListMap key field %q is missing", keyName),
		)
	}

	if keyMemberValue.IsNull() {
		return fieldpath.SelectorEntry{}, failure(
			keyPath,
			FailureNullKey,
			fmt.Sprintf("ListMap key field %q is null", keyName),
		)
	}

	selectorLiteral, err := literalFromValue(
		keyPath,
		keyMemberValue,
		keyFieldDescriptor.Type(),
		r.resolver,
		0,
	)
	if err != nil {
		return fieldpath.SelectorEntry{}, err
	}

	return fieldpath.NewSelectorEntry(keyName, selectorLiteral), nil
}
