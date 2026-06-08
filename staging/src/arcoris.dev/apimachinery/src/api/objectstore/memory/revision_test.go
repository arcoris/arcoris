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

package memory

import (
	"context"
	"testing"
)

func TestRevisionsAreMonotonicAcrossObjects(t *testing.T) {
	store := testStore(t)

	first := createState(t, store, testKey(1), "first")
	second := createState(t, store, testKey(2), "second")

	if !first.Revision.Before(second.Revision) {
		t.Fatalf("second revision %v did not advance from %v", second.Revision, first.Revision)
	}
}

func TestFailedGetDoesNotAdvanceRevision(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")

	_, ok, err := store.Get(context.Background(), testKey(2))
	requireNoError(t, err)
	if ok {
		t.Fatalf("missing object was found")
	}

	again, ok, err := store.Get(context.Background(), key)
	requireNoError(t, err)
	if !ok {
		t.Fatalf("created object missing")
	}
	if again.Revision != created.Revision {
		t.Fatalf("revision = %v; want %v", again.Revision, created.Revision)
	}
}
