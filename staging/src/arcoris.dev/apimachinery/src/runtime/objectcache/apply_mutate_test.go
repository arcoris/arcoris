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

package objectcache

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestCollectionApplyValidatedCreateUpdateDelete(t *testing.T) {
	key := testKey("system", 1)
	col := collection{}
	created := objectstore.MustCreatedChange(key, testState(key, 1, "created"))

	col.applyValidated(created)
	if col.revision != 1 || col.len() != 1 {
		t.Fatalf("after create revision=%s len=%d; want 1, 1", col.revision, col.len())
	}

	updated := objectstore.MustUpdatedChange(key, created.After, testState(key, 2, "updated"))
	col.applyValidated(updated)
	item, ok := col.item(key)
	if !ok || desiredString(t, item.State) != "updated" || col.revision != 2 {
		t.Fatalf("update did not replace item: item=%#v ok=%v revision=%s", item, ok, col.revision)
	}

	col.applyValidated(objectstore.MustDeletedChange(key, updated.After, 3))
	if col.revision != 3 || col.len() != 0 {
		t.Fatalf("after delete revision=%s len=%d; want 3, 0", col.revision, col.len())
	}
}
