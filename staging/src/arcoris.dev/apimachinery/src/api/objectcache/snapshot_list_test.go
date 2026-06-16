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

func TestSnapshotListQueriesMatchObjectQueryFullScan(t *testing.T) {
	source := testItems()
	snapshot := mustSnapshot(t, testListResult(14, source...))
	tests := []struct {
		name  string
		query objectquery.Query
	}{
		{
			name: "identity namespace",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustNamespaceEquals(t, "system"),
				},
			},
		},
		{
			name: "identity name",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Name: mustNameEquals(t, "worker-3"),
				},
			},
		},
		{
			name: "label exists",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelExists(t, "env")),
			},
		},
		{
			name: "label equals",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
			},
		},
		{
			name: "label in",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelIn(t, "env", "prod", "qa")),
			},
		},
		{
			name: "label does not exist",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelDoesNotExist(t, "env")),
			},
		},
		{
			name: "label not equals",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelNotEquals(t, "env", "prod")),
			},
		},
		{
			name: "label not in",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelNotIn(t, "env", "prod", "qa")),
			},
		},
		{
			name: "annotation exists",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationExists(t, "team")),
			},
		},
		{
			name: "annotation equals",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationEquals(t, "team", "core")),
			},
		},
		{
			name: "annotation in",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationIn(t, "team", "core", "tools")),
			},
		},
		{
			name: "annotation does not exist",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationDoesNotExist(t, "team")),
			},
		},
		{
			name: "annotation not equals",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationNotEquals(t, "team", "core")),
			},
		},
		{
			name: "annotation not in",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationNotIn(t, "team", "core", "tools")),
			},
		},
		{
			name: "combined",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustNamespaceEquals(t, "system"),
				},
				Labels: mustLabelSelector(
					t,
					mustLabelEquals(t, "tier", "backend"),
					mustLabelIn(t, "env", "prod", "qa"),
				),
				Annotations: mustAnnotationSelector(
					t,
					mustAnnotationEquals(t, "zone", "east"),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertSnapshotListMatchesObjectQueryFullScan(t, snapshot, source, tt.query)
		})
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

func TestSnapshotListInvalidQueryNote(t *testing.T) {
	// Malformed objectquery internals are tested inside api/objectquery.
	// Public objectquery constructors prevent constructing invalid queries here;
	// Snapshot.List delegates validation to objectquery.Compile.
	snapshot := mustSnapshot(t, testListResult(19, testItems()...))

	_, err := snapshot.List(objectquery.Query{})
	requireNoError(t, err)
}
