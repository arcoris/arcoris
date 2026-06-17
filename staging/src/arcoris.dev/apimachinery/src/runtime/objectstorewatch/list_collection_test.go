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
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestListCollectionReturnsValidCollectionRead(t *testing.T) {
	store := testRuntimeStore(t)
	key := testKey("system", 1)
	created := createObject(t, store, key, "created")

	read, err := store.ListCollection(context.Background(), testCollection())
	requireNoError(t, err)

	if read.Collection() != testCollection() {
		t.Fatalf("collection = %#v; want %#v", read.Collection(), testCollection())
	}
	if read.Revision() != created.Revision {
		t.Fatalf("revision = %s; want %s", read.Revision(), created.Revision)
	}
	if read.Len() != 1 {
		t.Fatalf("len = %d; want 1", read.Len())
	}
	if !read.IsValid() {
		t.Fatalf("collection read is invalid")
	}
}

func TestListCollectionAcceptsEmptyZeroRevisionResult(t *testing.T) {
	store := testRuntimeStore(t)

	read, err := store.ListCollection(context.Background(), testCollection())
	requireNoError(t, err)

	if read.Len() != 0 || !read.Revision().IsZero() {
		t.Fatalf("len/revision = %d/%s; want 0/0", read.Len(), read.Revision())
	}
}

func TestListCollectionRejectsInvalidCollection(t *testing.T) {
	store := testRuntimeStore(t)

	_, err := store.ListCollection(context.Background(), objectstore.ListRequest{})

	requireErrorIs(t, err, objectstore.ErrInvalidListRequest)
}

func TestListCollectionReturnsBackendErrors(t *testing.T) {
	store := testRuntimeStore(t)

	_, err := store.ListCollection(nil, testCollection())

	if err == nil || !errors.Is(err, objectstore.ErrNilContext) {
		t.Fatalf("err = %v; want nil context error", err)
	}
}
