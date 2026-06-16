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

func TestSnapshotListZeroQueryReturnsAllItems(t *testing.T) {
	source := testItems()
	snapshot := mustSnapshot(t, testListResult(14, source...))

	got, err := snapshot.List(objectquery.Query{})
	requireNoError(t, err)

	requireSameItems(t, got.Items, source)
	if got.Revision != 14 {
		t.Fatalf("Revision = %v; want 14", got.Revision)
	}
}

func TestSnapshotListRevisionPreservedForAllSomeAndNoMatches(t *testing.T) {
	source := testItems()
	snapshot := mustSnapshot(t, testListResult(17, source...))
	tests := []struct {
		name  string
		query objectquery.Query
	}{
		{name: "all", query: objectquery.Query{}},
		{
			name: "some",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
			},
		},
		{
			name: "none",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "missing")),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := snapshot.List(tt.query)
			requireNoError(t, err)
			if got.Revision != 17 {
				t.Fatalf("Revision = %v; want 17", got.Revision)
			}
		})
	}
}

func TestSnapshotListPreservesOrder(t *testing.T) {
	source := []objectstore.ListItem{
		testItem("system", "worker-3", 3, labelsMap("env", "prod"), nil),
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
		testItem("system", "worker-2", 2, labelsMap("env", "qa"), nil),
		testItem("system", "worker-4", 4, labelsMap("env", "prod"), nil),
	}
	snapshot := mustSnapshot(t, testListResult(18, source...))
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
	}

	got, err := snapshot.List(query)
	requireNoError(t, err)

	requireItemOrder(
		t,
		got.Items,
		itemRef{"system", "worker-3", 3},
		itemRef{"system", "worker-1", 1},
		itemRef{"system", "worker-4", 4},
	)
}

func TestSnapshotListInvalidQueryNote(t *testing.T) {
	// Malformed objectquery internals are tested inside api/objectquery.
	// Public objectquery constructors prevent constructing invalid queries here;
	// Snapshot.List delegates validation to objectquery.Compile.
	snapshot := mustSnapshot(t, testListResult(19, testItems()...))

	_, err := snapshot.List(objectquery.Query{})
	requireNoError(t, err)
}
