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

package objectstorewatch

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
)

func TestChangedEventClonesAndValidatesChange(t *testing.T) {
	store := testRuntimeStore(t)
	key := testKey("system", 1)
	created := createObject(t, store, key, "created")
	change := objectstore.Change{
		Kind:     objectstore.ChangeCreated,
		Key:      key,
		Revision: created.Revision,
		After:    created,
	}

	event, err := changedEvent(change)
	requireNoError(t, err)
	change.After.Revision = objectstore.Revision(99)

	if event.Kind != objectwatch.EventChanged {
		t.Fatalf("event kind = %s; want changed", event.Kind)
	}
	if event.Change.After.Revision != created.Revision {
		t.Fatalf("after revision = %s; want %s", event.Change.After.Revision, created.Revision)
	}
	requireNoError(t, event.Validate())
}

func TestChangedEventRejectsInvalidChange(t *testing.T) {
	_, err := changedEvent(objectstore.Change{})

	requireErrorIs(t, err, objectstore.ErrInvalidChange)
}
