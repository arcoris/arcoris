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

package objectmemorystore

import (
	"context"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestGetMissingReturnsFalseWithoutError(t *testing.T) {
	store := testStore(t)

	_, ok, err := store.Get(context.Background(), testKey(1))
	requireNoError(t, err)
	if ok {
		t.Fatalf("Get missing ok = true")
	}
}

func TestGetInvalidKeyReturnsInvalidKey(t *testing.T) {
	store := testStore(t)

	_, _, err := store.Get(context.Background(), objectstore.Key{})
	requireErrorIs(t, err, objectstore.ErrInvalidKey)
}

func TestGetReturnsCreatedState(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")

	got, ok, err := store.Get(context.Background(), key)
	requireNoError(t, err)
	if !ok {
		t.Fatalf("Get ok = false")
	}
	if got.Revision != created.Revision {
		t.Fatalf("Revision = %v; want %v", got.Revision, created.Revision)
	}
	requireDesiredString(t, got, "created")
}
