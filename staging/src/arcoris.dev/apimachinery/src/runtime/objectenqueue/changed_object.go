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

// ChangedObject returns a mapper that enqueues the object affected by a change.
//
// The mapper validates the committed change through objectstore.Change.Validate
// and then emits exactly one item for change.Key. It does not inspect payload
// fields, clone the change, or apply query semantics.
func ChangedObject() Mapper {
	return MapperFunc(mapChangedObject)
}

// mapChangedObject validates the committed change and emits its affected key.
func mapChangedObject(change objectstore.Change, emit EmitFunc) error {
	if err := change.Validate(); err != nil {
		return err
	}
	if emit == nil {
		return ErrNilEmit
	}

	return emit(objectworkqueue.Item{Key: change.Key})
}
