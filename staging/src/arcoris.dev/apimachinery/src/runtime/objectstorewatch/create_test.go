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

func TestCreatePublishesCreatedChange(t *testing.T) {
	store := testRuntimeStore(t)
	stream := watchAfter(t, store, testCollection(), 0)
	key := testKey("system", 1)

	created := createObject(t, store, key, "created")
	event := nextEvent(t, stream)

	requireChangedEvent(t, event, watchRequestAfter(t, testCollection(), 0), objectstore.ChangeCreated, created.Revision)
	if !event.Change.Key.Equal(key) {
		t.Fatalf("change key = %#v; want %#v", event.Change.Key, key)
	}
}
