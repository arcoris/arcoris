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

func TestCollectionValidateApplyAcceptsMatchingUpdate(t *testing.T) {
	key := testKey("system", 1)
	item := listItem(key, 1, "before")
	col, err := buildCollection(objectstore.ListResult{Items: []objectstore.ListItem{item}, Revision: 1})
	requireNoError(t, err)

	err = col.validateApply(objectstore.MustUpdatedChange(key, item.State, testState(key, 2, "after")))

	requireNoError(t, err)
}

func TestCollectionValidateApplyRejectsStateMismatch(t *testing.T) {
	key := testKey("system", 1)
	item := listItem(key, 1, "cached")
	col, err := buildCollection(objectstore.ListResult{Items: []objectstore.ListItem{item}, Revision: 1})
	requireNoError(t, err)

	err = col.validateApply(objectstore.MustDeletedChange(key, testState(key, 1, "different"), 2))

	requireErrorIs(t, err, ErrInvalidChange)
}
