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

package objectquery

import (
	"errors"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	"arcoris.dev/apimachinery/api/meta/annotations"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

// testResource is the canonical resource identity used by objectquery tests.
var testResource = apiidentity.GroupVersionResource{
	Group:    "control.arcoris.dev",
	Version:  "v1",
	Resource: "workers",
}

// TestQueryZeroAllNoneAndValidate locks down the public Query zero-value
// contract.
func TestQueryZeroAllNoneAndValidate(t *testing.T) {
	items := testItems()

	zero := mustPredicate(t, Query{})
	requireNames(t, zero.Filter(items), "worker-1", "worker-2", "worker-3", "worker-4")
	if !(Query{}).IsZero() {
		t.Fatal("zero Query IsZero() = false; want true")
	}
	requireNoError(t, (Query{}).Validate())

	all := mustPredicate(t, All())
	requireNames(t, all.Filter(items), "worker-1", "worker-2", "worker-3", "worker-4")

	none := mustPredicate(t, None())
	requireNames(t, none.Filter(items))
}

// TestQueryFromExprDetaches verifies private expression nodes are cloned before
// they become public Query values.
func TestQueryFromExprDetaches(t *testing.T) {
	source := termQuery(term{kind: termMetadata, metadataDomain: metadataLabels, metadataKey: "env", operator: OperatorEquals, stringValues: []string{"prod"}}).expr
	query := queryFromExpr(source)

	source.term.stringValues[0] = "qa"

	got := mustPredicate(t, query).Filter(testItems())
	requireNames(t, got, "worker-1", "worker-3")
}

// testItems returns a stable ordered collection used by predicate, plan, and
// projection tests.
func testItems() []objectstore.ListItem {
	return []objectstore.ListItem{
		testItem(
			"system",
			"worker-1",
			1,
			map[string]string{"env": "prod", "tier": "backend"},
			map[string]string{"team": "core"},
			desiredRecord("api", "prod", 3),
		),
		testItem(
			"system",
			"worker-2",
			2,
			map[string]string{"env": "qa", "tier": "backend"},
			map[string]string{"team": "tools"},
			desiredRecord("api", "qa", 1),
		),
		testItem(
			"other",
			"worker-3",
			3,
			map[string]string{"env": "prod", "tier": "frontend"},
			nil,
			desiredRecord("web", "prod", 5),
		),
		testItem("system", "worker-4", 4, nil, nil, value.MustRecordValue()),
	}
}

// testItem builds a committed list item whose key identity is authoritative for
// query matching.
func testItem(
	namespace metaidentity.Namespace,
	name metaidentity.Name,
	revision objectstore.Revision,
	rawLabels map[string]string,
	rawAnnotations map[string]string,
	desired value.Value,
) objectstore.ListItem {
	objectName := metaidentity.ObjectName{Namespace: namespace, Name: name}
	objectMeta := meta.ObjectMeta{
		Name:        name,
		Namespace:   namespace,
		Labels:      mustLabels(rawLabels),
		Annotations: mustAnnotations(rawAnnotations),
	}
	obj := object.New[value.Value, value.Value](
		meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
			Group:   testResource.Group,
			Version: testResource.Version,
			Kind:    "Worker",
		}),
		objectMeta,
		desired,
	)

	return objectstore.ListItem{
		Key: objectstore.MustKey(testResource, objectName),
		State: objectstore.State{
			Object:    obj,
			Ownership: objectownership.EmptyState(),
			Revision:  revision,
		},
	}
}

// desiredRecord builds a small nested desired payload for selectable field
// tests.
func desiredRecord(image string, phase string, replicas int64) value.Value {
	return value.MustRecordValue(
		value.MustRecordMember("spec", value.MustRecordValue(
			value.MustRecordMember("image", value.StringValue(image)),
			value.MustRecordMember("phase", value.StringValue(phase)),
			value.MustRecordMember("replicas", value.Int64Value(replicas)),
			value.MustRecordMember("nullable", value.NullValue()),
		)),
	)
}

// mustLabels converts string fixtures into typed labels and panics on invalid
// test data.
func mustLabels(values map[string]string) labels.Set {
	set, err := labels.FromStrings(values)
	if err != nil {
		panic(err)
	}

	return set
}

// mustAnnotations converts string fixtures into typed annotations and panics on
// invalid test data.
func mustAnnotations(values map[string]string) annotations.Set {
	set, err := annotations.FromStrings(values)
	if err != nil {
		panic(err)
	}

	return set
}

// mustQ unwraps constructor results in tests where invalid input would be a
// fixture bug.
func mustQ(query Query, err error) Query {
	if err != nil {
		panic(err)
	}

	return query
}

// mustPredicate compiles a query for tests and fails immediately on errors.
func mustPredicate(t testing.TB, query Query, opts ...CompileOption) Predicate {
	t.Helper()

	predicate, err := Compile(query, opts...)
	requireNoError(t, err)
	return predicate
}

// mustAnd builds an AND query from already-valid test queries.
func mustAnd(t testing.TB, queries ...Query) Query {
	t.Helper()
	return mustQ(And(queries...))
}

// mustOr builds an OR query from already-valid test queries.
func mustOr(t testing.TB, queries ...Query) Query {
	t.Helper()
	return mustQ(Or(queries...))
}

// mustNot builds a NOT query from an already-valid test query.
func mustNot(t testing.TB, query Query) Query {
	t.Helper()
	return mustQ(Not(query))
}

// requireNoError fails a test when err is non-nil.
func requireNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// requireErrorIs asserts errors.Is without obscuring the original error text.
func requireErrorIs(t testing.TB, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

// requireNames asserts object names and order for filtered list results.
func requireNames(t *testing.T, items []objectstore.ListItem, names ...metaidentity.Name) {
	t.Helper()
	if len(items) != len(names) {
		t.Fatalf("len = %d; want %d", len(items), len(names))
	}
	for i, name := range names {
		if items[i].Key.Object.Name != name {
			t.Fatalf("items[%d].name = %s; want %s", i, items[i].Key.Object.Name, name)
		}
	}
}

// requireNamesFromStrings is a small adapter for table tests whose expected
// names are easier to read as string literals.
func requireNamesFromStrings(t *testing.T, items []objectstore.ListItem, names ...string) {
	t.Helper()

	typed := make([]metaidentity.Name, len(names))
	for i, name := range names {
		typed[i] = metaidentity.Name(name)
	}
	requireNames(t, items, typed...)
}
