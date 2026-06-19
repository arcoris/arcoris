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
)

func TestCreatedChangeBuildsValidDetachedChange(t *testing.T) {
	store := testRuntimeStore(t)
	key := testKey("system", 1)
	created := createObject(t, store, key, "created")

	change, err := createdChange(key, created)
	requireNoError(t, err)
	created.Revision = 99

	if change.Kind != objectstore.ChangeCreated {
		t.Fatalf("kind = %s; want created", change.Kind)
	}
	if change.After.Revision == created.Revision {
		t.Fatalf("change retained caller mutation")
	}
	requireNoError(t, change.Validate())
}

func TestUpdatedChangeBuildsValidDetachedChange(t *testing.T) {
	store := testRuntimeStore(t)
	key := testKey("system", 1)
	before := createObject(t, store, key, "created")
	after := updateObject(t, store, key, before.Revision, "updated")

	change, err := updatedChange(key, before, after)
	requireNoError(t, err)
	after.Revision = 99

	if change.Kind != objectstore.ChangeUpdated {
		t.Fatalf("kind = %s; want updated", change.Kind)
	}
	if change.After.Revision == after.Revision {
		t.Fatalf("change retained caller mutation")
	}
	requireNoError(t, change.Validate())
}

func TestDeletedChangeUsesDeleteRevision(t *testing.T) {
	store := testRuntimeStore(t)
	key := testKey("system", 1)
	created := createObject(t, store, key, "created")
	deleted := deleteObject(t, store, key, created.Revision)

	change, err := deletedChange(key, deleted)
	requireNoError(t, err)

	if change.Kind != objectstore.ChangeDeleted {
		t.Fatalf("kind = %s; want deleted", change.Kind)
	}
	if change.Revision != deleted.Revision {
		t.Fatalf("revision = %s; want delete revision %s", change.Revision, deleted.Revision)
	}
	if change.Before.Revision != created.Revision {
		t.Fatalf("before revision = %s; want live revision %s", change.Before.Revision, created.Revision)
	}
	requireNoError(t, change.Validate())
}

func TestChangeHelpersRejectInvalidChanges(t *testing.T) {
	if _, err := createdChange(objectstore.Key{}, objectstore.State{}); err == nil {
		t.Fatalf("createdChange() error = nil; want error")
	}
	if _, err := updatedChange(objectstore.Key{}, objectstore.State{}, objectstore.State{}); err == nil {
		t.Fatalf("updatedChange() error = nil; want error")
	}
	if _, err := deletedChange(objectstore.Key{}, objectstore.DeleteResult{}); err == nil {
		t.Fatalf("deletedChange() error = nil; want error")
	}
}
