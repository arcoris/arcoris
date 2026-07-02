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

// Mapper maps one committed object change to zero or more reconciliation items.
type Mapper interface {
	Map(objectstore.Change, EmitFunc) error
}

// MapperFunc adapts a function to Mapper.
type MapperFunc func(objectstore.Change, EmitFunc) error

// Map calls f with change and emit.
func (f MapperFunc) Map(change objectstore.Change, emit EmitFunc) error {
	if f == nil {
		return ErrNilMapper
	}

	return f(change, emit)
}

// ChangedObject returns a mapper that enqueues the object affected by a change.
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
