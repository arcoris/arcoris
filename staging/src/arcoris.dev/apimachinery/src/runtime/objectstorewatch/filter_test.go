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

func TestWatchFiltersByNamespaceAndResource(t *testing.T) {
	store := testRuntimeStore(t)
	stream := watchAfter(t, store, namespaceCollection("alpha"), 0)

	alpha := createObject(t, store, testKey("alpha", 1), "alpha")
	createObject(t, store, testKey("beta", 1), "beta")
	_, err := store.Create(context.Background(), otherResourceKey("alpha", 1), stateForKey(otherResourceKey("alpha", 1), "other"))
	requireNoError(t, err)

	event := nextEvent(t, stream)
	requireChangedEvent(t, event, watchRequestAfter(t, namespaceCollection("alpha"), 0), objectstore.ChangeCreated, alpha.Revision)
	requireNoEvent(t, stream)
}

func TestWatchAllNamespacesReceivesMultipleNamespaces(t *testing.T) {
	store := testRuntimeStore(t)
	stream := watchAfter(t, store, testCollection(), 0)

	first := createObject(t, store, testKey("alpha", 1), "alpha")
	second := createObject(t, store, testKey("beta", 1), "beta")

	if event := nextEvent(t, stream); event.Revision != first.Revision {
		t.Fatalf("first revision = %s; want %s", event.Revision, first.Revision)
	}
	if event := nextEvent(t, stream); event.Revision != second.Revision {
		t.Fatalf("second revision = %s; want %s", event.Revision, second.Revision)
	}
}
