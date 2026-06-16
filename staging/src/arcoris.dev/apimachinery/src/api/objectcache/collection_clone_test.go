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

	"arcoris.dev/apimachinery/api/objectquery"
)

func TestCollectionCloneDetachesItemsAndIndexes(t *testing.T) {
	col := mustCollection(t, testListResult(
		7,
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	))

	cloned := col.clone()
	mutated := cloned.items[cloned.order[0]]
	mutateItem(&mutated, "mutated")
	cloned.items[cloned.order[0]] = mutated
	cloned.indexes.remove(mutated)

	got, ok := col.item(col.order[0])
	if !ok {
		t.Fatal("item() ok = false; want true")
	}
	if got := desiredString(t, got); got != "worker-1" {
		t.Fatalf("Desired = %q; want worker-1", got)
	}
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
	}
	assertCollectionListEquivalent(t, col, query)
}
