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

func TestUpdatePublishesUpdatedChangeWithBeforeAndAfter(t *testing.T) {
	store := testRuntimeStore(t)
	key := testKey("system", 1)
	created := createObject(t, store, key, "created")
	stream := watchAfter(t, store, testCollection(), created.Revision)

	updated := updateObject(t, store, key, created.Revision, "updated")
	event := nextEvent(t, stream)

	requireChangedEvent(t, event, watchRequestAfter(t, testCollection(), created.Revision), objectstore.ChangeUpdated, updated.Revision)
	requireDesiredString(t, event.Change.Before, "created")
	requireDesiredString(t, event.Change.After, "updated")
	if !event.Change.Before.Revision.Before(event.Change.After.Revision) {
		t.Fatalf("before revision %s is not before after revision %s", event.Change.Before.Revision, event.Change.After.Revision)
	}
}
