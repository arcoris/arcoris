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
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestChangedObjectNilEmit(t *testing.T) {
	err := ChangedObject().Map(createdChange(t, 1), nil)
	requireErrorIs(t, err, ErrNilEmit)
}

func TestChangedObjectRejectsInvalidChange(t *testing.T) {
	var emitted bool
	err := ChangedObject().Map(objectstore.Change{}, func(objectworkqueue.Item) error {
		emitted = true
		return nil
	})

	requireErrorIs(t, err, objectstore.ErrInvalidChange)
	if emitted {
		t.Fatalf("emitted item for invalid change")
	}
}

func TestChangedObjectEmitsChangedKey(t *testing.T) {
	tests := []struct {
		name   string
		change objectstore.Change
	}{
		{name: "created", change: createdChange(t, 1)},
		{name: "updated", change: updatedChange(t, 2)},
		{name: "deleted", change: deletedChange(t, 3)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []objectworkqueue.Item
			err := ChangedObject().Map(tt.change, func(item objectworkqueue.Item) error {
				got = append(got, item)
				return nil
			})

			requireNoError(t, err)
			if len(got) != 1 {
				t.Fatalf("items = %d; want 1", len(got))
			}
			requireItem(t, got[0], objectworkqueue.Item{Key: tt.change.Key})
		})
	}
}

func TestChangedObjectPropagatesEmitError(t *testing.T) {
	wantErr := errors.New("emit failed")
	err := ChangedObject().Map(createdChange(t, 1), func(objectworkqueue.Item) error {
		return wantErr
	})

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
}

func TestMapperFuncNil(t *testing.T) {
	var mapper MapperFunc
	err := mapper.Map(createdChange(t, 1), func(objectworkqueue.Item) error {
		return nil
	})

	requireErrorIs(t, err, ErrNilMapper)
}

func TestMapperFuncDelegates(t *testing.T) {
	change := createdChange(t, 1)
	wantErr := errors.New("delegate failed")
	called := false
	mapper := MapperFunc(func(got objectstore.Change, emit EmitFunc) error {
		called = true
		requireChange(t, got, change)
		if emit == nil {
			t.Fatalf("emit is nil")
		}
		return wantErr
	})

	err := mapper.Map(change, func(objectworkqueue.Item) error { return nil })

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
	if !called {
		t.Fatalf("mapper function was not called")
	}
}
