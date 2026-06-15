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

package objectindex

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestSelectZeroQueryReturnsAllItemsInInputOrder(t *testing.T) {
	items := testItems()

	got, err := Build(items).Select(objectquery.Query{})
	requireNoError(t, err)

	requireSameItems(t, got, items)
}

func TestSelectPredicateZeroReturnsAllItemsInInputOrder(t *testing.T) {
	items := testItems()

	got := Build(items).SelectPredicate(objectquery.Predicate{})

	requireSameItems(t, got, items)
}

func TestSelectNoMatchReturnsNil(t *testing.T) {
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "missing")),
	}

	got, err := Build(testItems()).Select(query)
	requireNoError(t, err)

	if got != nil {
		t.Fatalf("Select() = %#v; want nil", got)
	}
}

func TestSelectIdentityQueriesMatchFullScan(t *testing.T) {
	items := testItems()
	tests := []struct {
		name  string
		query objectquery.Query
	}{
		{
			name: "namespace",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustNamespaceEquals(t, "system"),
				},
			},
		},
		{
			name: "name",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Name: mustNameEquals(t, "worker-3"),
				},
			},
		},
		{
			name: "namespace and name",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustNamespaceEquals(t, "system"),
					Name:      mustNameEquals(t, "worker-2"),
				},
			},
		},
		{
			name: "no match",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustNamespaceEquals(t, "missing"),
					Name:      mustNameEquals(t, "worker-2"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertIndexMatchesFullScan(t, items, tt.query)
		})
	}
}

func TestSelectLabelQueriesMatchFullScan(t *testing.T) {
	items := testItems()
	tests := []struct {
		name  string
		query objectquery.Query
	}{
		{
			name: "exists",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelExists(t, "env")),
			},
		},
		{
			name: "equals",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
			},
		},
		{
			name: "in",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelIn(t, "env", "qa", "prod")),
			},
		},
		{
			name: "does not exist",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelDoesNotExist(t, "env")),
			},
		},
		{
			name: "not equals",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelNotEquals(t, "env", "prod")),
			},
		},
		{
			name: "not in",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelNotIn(t, "env", "prod", "qa")),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertIndexMatchesFullScan(t, items, tt.query)
		})
	}
}

func TestSelectAnnotationQueriesMatchFullScan(t *testing.T) {
	items := testItems()
	tests := []struct {
		name  string
		query objectquery.Query
	}{
		{
			name: "exists",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationExists(t, "team")),
			},
		},
		{
			name: "equals",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationEquals(t, "team", "core")),
			},
		},
		{
			name: "in",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationIn(t, "team", "core", "tools")),
			},
		},
		{
			name: "does not exist",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationDoesNotExist(t, "team")),
			},
		},
		{
			name: "not equals",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationNotEquals(t, "team", "core")),
			},
		},
		{
			name: "not in",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationNotIn(t, "team", "core", "tools")),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertIndexMatchesFullScan(t, items, tt.query)
		})
	}
}

func TestSelectCombinedIdentityLabelsAnnotationsMatchesFullScan(t *testing.T) {
	items := testItems()
	query := objectquery.Query{
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
	}

	assertIndexMatchesFullScan(t, items, query)
}

func TestSelectMultipleLabelRequirementsUseAndSemantics(t *testing.T) {
	items := testItems()
	query := objectquery.Query{
		Labels: mustLabelSelector(
			t,
			mustLabelEquals(t, "env", "prod"),
			mustLabelEquals(t, "tier", "backend"),
		),
	}

	assertIndexMatchesFullScan(t, items, query)
}

func TestSelectMultipleAnnotationRequirementsUseAndSemantics(t *testing.T) {
	items := testItems()
	query := objectquery.Query{
		Annotations: mustAnnotationSelector(
			t,
			mustAnnotationEquals(t, "team", "core"),
			mustAnnotationEquals(t, "zone", "west"),
		),
	}

	assertIndexMatchesFullScan(t, items, query)
}

func TestSelectLabelInUsesUnionWithinRequirementAndAndAcrossRequirements(t *testing.T) {
	items := testItems()
	query := objectquery.Query{
		Labels: mustLabelSelector(
			t,
			mustLabelIn(t, "env", "prod", "qa"),
			mustLabelEquals(t, "tier", "backend"),
		),
	}

	assertIndexMatchesFullScan(t, items, query)
}

func TestSelectNegativeOnlyQueryFallsBackToFullScanSemantics(t *testing.T) {
	items := testItems()
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelNotEquals(t, "env", "prod")),
	}

	assertIndexMatchesFullScan(t, items, query)
}

func TestSelectPositiveAndNegativeQueryUsesResidualPredicate(t *testing.T) {
	items := testItems()
	query := objectquery.Query{
		Labels: mustLabelSelector(
			t,
			mustLabelExists(t, "env"),
			mustLabelNotEquals(t, "tier", "frontend"),
		),
	}

	assertIndexMatchesFullScan(t, items, query)
}

func TestSelectAbsentKeyNegativeSemanticsMatchObjectQuery(t *testing.T) {
	items := testItems()
	query := objectquery.Query{
		Labels: mustLabelSelector(
			t,
			mustLabelNotIn(t, "missing", "a", "b"),
		),
		Annotations: mustAnnotationSelector(
			t,
			mustAnnotationNotEquals(t, "missing", "value"),
		),
	}

	assertIndexMatchesFullScan(t, items, query)
}

func TestSelectPreservesInputOrderForMultipleMatches(t *testing.T) {
	items := []objectstore.ListItem{
		testItem("system", "worker-3", 3, labelsMap("env", "prod"), nil),
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
		testItem("system", "worker-2", 2, labelsMap("env", "qa"), nil),
		testItem("system", "worker-4", 4, labelsMap("env", "prod"), nil),
	}
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
	}

	got, err := Build(items).Select(query)
	requireNoError(t, err)

	requireItemOrder(
		t,
		got,
		itemRef{"system", "worker-3", 3},
		itemRef{"system", "worker-1", 1},
		itemRef{"system", "worker-4", 4},
	)
}

func TestSelectPreservesInputOrderAfterInUnion(t *testing.T) {
	items := testItems()
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelIn(t, "env", "qa", "prod")),
	}

	got, err := Build(items).Select(query)
	requireNoError(t, err)

	requireItemOrder(
		t,
		got,
		itemRef{"system", "worker-1", 1},
		itemRef{"system", "worker-2", 2},
		itemRef{"other", "worker-3", 3},
		itemRef{"", "cluster-worker", 4},
	)
}

func TestSelectPreservesDuplicateKeyMultiplicityAndOrder(t *testing.T) {
	items := []objectstore.ListItem{
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
		testItem("system", "worker-1", 2, labelsMap("env", "qa"), nil),
		testItem("system", "worker-1", 3, labelsMap("env", "prod"), nil),
	}
	query := objectquery.Query{
		Identity: objectquery.IdentitySelector{
			Name: mustNameEquals(t, "worker-1"),
		},
	}

	got, err := Build(items).Select(query)
	requireNoError(t, err)

	requireItemOrder(
		t,
		got,
		itemRef{"system", "worker-1", 1},
		itemRef{"system", "worker-1", 2},
		itemRef{"system", "worker-1", 3},
	)
}
