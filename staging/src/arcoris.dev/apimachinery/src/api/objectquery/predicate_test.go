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
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	"arcoris.dev/apimachinery/api/meta/annotations"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

func TestPredicateQueryZeroPredicate(t *testing.T) {
	var predicate Predicate
	query := predicate.Query()

	if !query.Identity.IsZero() || !query.Labels.IsZero() || !query.Annotations.IsZero() {
		t.Fatalf("Query() = %#v; want zero query", query)
	}
}

func TestPredicateQueryReturnsCanonicalQuery(t *testing.T) {
	predicate, err := Compile(Query{
		Labels: mustLabelSelector(
			t,
			mustLabelEquals(t, "tier", "backend"),
			mustLabelIn(t, "env", "qa", "prod"),
		),
	})
	requireNoError(t, err)

	requirements := predicate.Query().Labels.Requirements()
	if len(requirements) != 2 {
		t.Fatalf("len = %d; want 2", len(requirements))
	}
	if requirements[0].Key() != "env" || requirements[1].Key() != "tier" {
		t.Fatalf("requirement order = %q, %q; want env, tier", requirements[0].Key(), requirements[1].Key())
	}
}

func TestPredicateQueryPreservesIdentity(t *testing.T) {
	predicate, err := Compile(Query{Identity: mustWithObject(t, "system", "worker")})
	requireNoError(t, err)

	query := predicate.Query()
	namespace, ok := query.Identity.Namespace.Namespace()
	if !ok || namespace != "system" {
		t.Fatalf("Namespace() = %q, %v; want system, true", namespace, ok)
	}
	name, ok := query.Identity.Name.Name()
	if !ok || name != "worker" {
		t.Fatalf("Name() = %q, %v; want worker, true", name, ok)
	}
}

func TestPredicateQueryDefensiveCopy(t *testing.T) {
	predicate, err := Compile(Query{
		Labels: mustLabelSelector(t, mustLabelIn(t, "env", "qa", "prod")),
	})
	requireNoError(t, err)

	query := predicate.Query()
	values := query.Labels.Requirements()[0].Values()
	values[0] = "mutated"
	query.Labels.requirements[0].req.values[0] = "mutated"
	query.Labels.requirements = append(query.Labels.requirements, mustLabelEquals(t, "tier", "backend"))

	next := predicate.Query().Labels.Requirements()
	if len(next) != 1 {
		t.Fatalf("len = %d; want 1", len(next))
	}
	requireStrings(t, next[0].Values(), "prod", "qa")
	if !predicate.Match(testItem("system", "worker", map[string]string{"env": "prod"}, nil)) {
		t.Fatal("predicate was mutated through Query()")
	}
}

func testItem(namespace string, name string, labelValues map[string]string, annotationValues map[string]string) objectstore.ListItem {
	labelSet, err := labels.FromStrings(labelValues)
	if err != nil {
		panic(err)
	}
	annotationSet, err := annotations.FromStrings(annotationValues)
	if err != nil {
		panic(err)
	}

	key := objectstore.MustKey(
		apiidentity.GroupVersionResource{
			Group:    "control.arcoris.dev",
			Version:  "v1",
			Resource: "workers",
		},
		metaidentity.ObjectName{
			Namespace: metaidentity.Namespace(namespace),
			Name:      metaidentity.Name(name),
		},
	)

	return objectstore.ListItem{
		Key: key,
		State: objectstore.State{
			Object: object.New[value.Value, value.Value](
				meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
					Group:   "control.arcoris.dev",
					Version: "v1",
					Kind:    "Worker",
				}),
				meta.ObjectMeta{
					Name:        metaidentity.Name(name),
					Namespace:   metaidentity.Namespace(namespace),
					Labels:      labelSet,
					Annotations: annotationSet,
				},
				value.StringValue(name),
			),
			Revision: 1,
		},
	}
}
