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

func TestCollectionApplyUnknownKind(t *testing.T) {
	col := mustCollection(t, testListResult(1))

	err := col.validateApply(objectstore.Change{Kind: objectstore.ChangeKind(99)})

	requireErrorIs(t, err, ErrInvalidChange)
}

func TestCollectionValidateApplyDoesNotMutateState(t *testing.T) {
	item := testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil)
	col := mustCollection(t, testListResult(10, item))
	before := col.listItems()
	change := objectstore.MustCreatedChange(item.Key, testItem(
		"system",
		"worker-1",
		11,
		labelsMap("env", "qa"),
		nil,
	).State)

	err := col.validateApply(change)

	requireErrorIs(t, err, ErrInvalidChange)
	if col.revision != 10 {
		t.Fatalf("revision = %v; want 10", col.revision)
	}
	requireSameItems(t, col.listItems(), before)
	assertCollectionInvariants(t, col)
}

func TestCollectionApplyValidatedMutatesIncrementally(t *testing.T) {
	item := testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil)
	col := mustCollection(t, testListResult(1, item))
	next := testItem("system", "worker-1", 2, labelsMap("env", "qa"), nil)
	change := objectstore.MustUpdatedChange(item.Key, item.State, next.State)

	requireNoError(t, col.validateApply(change))
	col.applyValidated(change.Clone())

	requireItemOrder(t, col.listItems(), itemRef{"system", "worker-1", 2})
	assertCollectionInvariants(t, col)
}

func TestCollectionRemoveMissingOrderKeyLeavesOrder(t *testing.T) {
	item := testItem("system", "worker-1", 1, nil, nil)
	col := mustCollection(t, testListResult(1, item))

	col.removeOrderKey(testItem("system", "missing", 2, nil, nil).Key)

	requireItemOrder(t, col.listItems(), itemRef{"system", "worker-1", 1})
}
