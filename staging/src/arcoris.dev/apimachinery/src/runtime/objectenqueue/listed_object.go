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

package objectenqueue

import (
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

// ListedObject returns a mapper that enqueues the object carried by a list item.
//
// The mapper validates the live list item through objectstore.ValidateListItem
// and then emits exactly one item for item.Key. It does not inspect payload
// fields, clone the item, retain the item, or apply query semantics.
func ListedObject() ListItemMapper {
	return ListItemMapperFunc(mapListedObject)
}

// mapListedObject validates one live listed object and emits its key.
func mapListedObject(item objectstore.ListItem, emit EmitFunc) error {
	if err := objectstore.ValidateListItem(item); err != nil {
		return err
	}
	if emit == nil {
		return ErrNilEmit
	}

	return emit(objectworkqueue.Item{Key: item.Key})
}
