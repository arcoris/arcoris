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
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestSnapshotListReturnsDetachedCopies(t *testing.T) {
	source := []objectstore.ListItem{
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	}
	snapshot := mustSnapshot(t, testListResult(19, source...))
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
	}

	got, err := snapshot.List(query)
	requireNoError(t, err)
	mutateItem(&got.Items[0], "mutated")

	again, err := snapshot.List(query)
	requireNoError(t, err)
	if got := desiredString(t, again.Items[0]); got != "worker-1" {
		t.Fatalf("Desired = %q; want worker-1", got)
	}
	if got := labelValue(t, again.Items[0], "env"); got != "prod" {
		t.Fatalf("label env = %q; want prod", got)
	}
}

func TestSnapshotUnaffectedByReturnedListMutation(t *testing.T) {
	source := []objectstore.ListItem{
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	}
	snapshot := mustSnapshot(t, testListResult(19, source...))

	got, err := snapshot.List(objectquery.Query{})
	requireNoError(t, err)
	mutateItem(&got.Items[0], "changed")
	got.Items = append(got.Items, testItem("system", "extra", 99, nil, nil))

	again := snapshot.Items()
	requireItemOrder(t, again, itemRef{"system", "worker-1", 1})
	if got := desiredString(t, again[0]); got != "worker-1" {
		t.Fatalf("Desired = %q; want worker-1", got)
	}
}
