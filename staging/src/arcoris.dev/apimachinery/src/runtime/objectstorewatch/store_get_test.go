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
)

func TestGetDelegatesToBackend(t *testing.T) {
	store := testRuntimeStore(t)
	key := testKey("system", 1)
	created := createObject(t, store, key, "created")

	got, ok, err := store.Get(context.Background(), key)
	requireNoError(t, err)

	if !ok {
		t.Fatalf("ok = false; want true")
	}
	if got.Revision != created.Revision {
		t.Fatalf("revision = %s; want %s", got.Revision, created.Revision)
	}
	requireDesiredString(t, got, "created")
}

func TestGetMissingDelegatesToBackend(t *testing.T) {
	store := testRuntimeStore(t)

	got, ok, err := store.Get(context.Background(), testKey("system", 1))
	requireNoError(t, err)

	if ok {
		t.Fatalf("ok = true; want false")
	}
	if !got.Revision.IsZero() {
		t.Fatalf("revision = %s; want zero", got.Revision)
	}
}
