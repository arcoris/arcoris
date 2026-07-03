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

func TestListedObjectNilEmit(t *testing.T) {
	err := ListedObject().MapListItem(testListItem(1, 1, "listed"), nil)
	requireErrorIs(t, err, ErrNilEmit)
}

func TestListedObjectRejectsInvalidListItem(t *testing.T) {
	var emitted bool
	err := ListedObject().MapListItem(objectstore.ListItem{}, func(objectworkqueue.Item) error {
		emitted = true
		return nil
	})

	requireErrorIs(t, err, objectstore.ErrInvalidListResult)
	if emitted {
		t.Fatalf("emitted item for invalid list item")
	}
}

func TestListedObjectEmitsListedKey(t *testing.T) {
	item := testListItem(1, 1, "listed")
	var got []objectworkqueue.Item

	err := ListedObject().MapListItem(item, func(item objectworkqueue.Item) error {
		got = append(got, item)
		return nil
	})

	requireNoError(t, err)
	if len(got) != 1 {
		t.Fatalf("items = %d; want 1", len(got))
	}
	requireItem(t, got[0], objectworkqueue.Item{Key: item.Key})
}

func TestListedObjectPropagatesEmitError(t *testing.T) {
	wantErr := errors.New("emit failed")
	err := ListedObject().MapListItem(testListItem(1, 1, "listed"), func(objectworkqueue.Item) error {
		return wantErr
	})

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
}
