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
	"context"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestListDelegatesToBackend(t *testing.T) {
	store := testRuntimeStore(t)
	key := testKey("system", 1)
	created := createObject(t, store, key, "created")

	result, err := store.List(context.Background(), namespaceCollection("system"))
	requireNoError(t, err)

	if result.Len() != 1 {
		t.Fatalf("len = %d; want 1", result.Len())
	}
	if result.Revision != created.Revision {
		t.Fatalf("revision = %s; want %s", result.Revision, created.Revision)
	}
	if result.Items[0].Key != key {
		t.Fatalf("key = %#v; want %#v", result.Items[0].Key, key)
	}
	requireDesiredString(t, result.Items[0].State, "created")
}

func TestListPreservesBackendValidation(t *testing.T) {
	store := testRuntimeStore(t)

	_, err := store.List(context.Background(), objectstore.ListRequest{})

	requireErrorIs(t, err, objectstore.ErrInvalidListRequest)
}
