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

package objectlifecycle

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/meta/annotations"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

func TestListZeroQueryPreservesExistingListBehavior(t *testing.T) {
	executor := testExecutor(t)
	first := createObject(t, executor, 1, "api:v1", owner("creator"))
	second := createObject(t, executor, 2, "api:v2", owner("creator"))

	result, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
		Query:    objectquery.Query{},
	})
	requireNoError(t, err)

	if result.Len() != 2 {
		t.Fatalf("len = %d; want 2", result.Len())
	}
	requireLifecycleListItem(
		t,
		result.Items[0],
		objectstore.MustKey(testGVR(), testName(1)),
		first.State.Revision,
		"api:v1",
	)
	requireLifecycleListItem(
		t,
		result.Items[1],
		objectstore.MustKey(testGVR(), testName(2)),
		second.State.Revision,
		"api:v2",
	)
	if result.Revision != second.State.Revision {
		t.Fatalf("revision = %v; want %v", result.Revision, second.State.Revision)
	}
}

func TestListLabelQueryFiltersResult(t *testing.T) {
	executor := testExecutor(t)
	createObjectWithMetadata(
		t,
		executor,
		1,
		"system",
		"api:prod",
		map[string]string{"env": "prod"},
		nil,
	)
	createObjectWithMetadata(
		t,
		executor,
		2,
		"system",
		"api:qa",
		map[string]string{"env": "qa"},
		nil,
	)
	latest := createObjectWithMetadata(t, executor, 3, "system", "api:none", nil, nil)

	result, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
		Query: objectquery.Query{
			Labels: mustLifecycleLabelSelector(t, mustLifecycleLabelEquals(t, "env", "prod")),
		},
	})
	requireNoError(t, err)

	requireLifecycleListNames(t, result, "worker-1")
	if result.Revision != latest.State.Revision {
		t.Fatalf("revision = %v; want %v", result.Revision, latest.State.Revision)
	}
	requireImage(t, result.Items[0].State, "api:prod")
}

func TestListAnnotationQueryFiltersResult(t *testing.T) {
	executor := testExecutor(t)
	createObjectWithMetadata(
		t,
		executor,
		1,
		"system",
		"api:backend",
		nil,
		map[string]string{"tier": "backend"},
	)
	createObjectWithMetadata(
		t,
		executor,
		2,
		"system",
		"api:frontend",
		nil,
		map[string]string{"tier": "frontend"},
	)

	result, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
		Query: objectquery.Query{
			Annotations: mustLifecycleAnnotationSelector(
				t,
				mustLifecycleAnnotationEquals(t, "tier", "backend"),
			),
		},
	})
	requireNoError(t, err)

	requireLifecycleListNames(t, result, "worker-1")
	requireImage(t, result.Items[0].State, "api:backend")
}

func TestListIdentityNameQueryFiltersResult(t *testing.T) {
	executor := testExecutor(t)
	createObject(t, executor, 1, "api:v1", owner("creator"))
	createObject(t, executor, 2, "api:v2", owner("creator"))

	result, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
		Query: objectquery.Query{
			Identity: objectquery.IdentitySelector{
				Name: mustLifecycleNameEquals(t, "worker-2"),
			},
		},
	})
	requireNoError(t, err)

	requireLifecycleListNames(t, result, "worker-2")
	requireImage(t, result.Items[0].State, "api:v2")
}

func TestListCombinedQueryUsesAndSemantics(t *testing.T) {
	executor := testExecutor(t)
	createObjectWithMetadata(
		t,
		executor,
		1,
		"system",
		"api:match",
		map[string]string{"env": "prod"},
		map[string]string{"tier": "backend"},
	)
	createObjectWithMetadata(
		t,
		executor,
		2,
		"system",
		"api:wrong-name",
		map[string]string{"env": "prod"},
		map[string]string{"tier": "backend"},
	)
	createObjectWithMetadata(
		t,
		executor,
		3,
		"system",
		"api:wrong-label",
		map[string]string{"env": "qa"},
		map[string]string{"tier": "backend"},
	)
	createObjectWithMetadata(
		t,
		executor,
		4,
		"system",
		"api:wrong-annotation",
		map[string]string{"env": "prod"},
		map[string]string{"tier": "frontend"},
	)

	result, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
		Query: objectquery.Query{
			Identity: objectquery.IdentitySelector{
				Name: mustLifecycleNameEquals(t, "worker-1"),
			},
			Labels: mustLifecycleLabelSelector(t, mustLifecycleLabelEquals(t, "env", "prod")),
			Annotations: mustLifecycleAnnotationSelector(
				t,
				mustLifecycleAnnotationEquals(t, "tier", "backend"),
			),
		},
	})
	requireNoError(t, err)

	requireLifecycleListNames(t, result, "worker-1")
	requireImage(t, result.Items[0].State, "api:match")
}

func TestListQueryNamespacePolicyMatrix(t *testing.T) {
	tests := []struct {
		name            string
		globalResource  bool
		scope           func(t *testing.T) objectstore.ListScope
		queryNamespace  metaidentity.Namespace
		hasNamespace    bool
		wantReason      ErrorReason
		wantStoreCalled bool
	}{
		{
			name:            "namespace scope different query namespace rejected",
			scope:           lifecycleNamespaceScope("alpha"),
			queryNamespace:  "beta",
			hasNamespace:    true,
			wantReason:      ErrorReasonInvalidQueryScope,
			wantStoreCalled: false,
		},
		{
			name:            "namespace scope same query namespace allowed",
			scope:           lifecycleNamespaceScope("alpha"),
			queryNamespace:  "alpha",
			hasNamespace:    true,
			wantStoreCalled: true,
		},
		{
			name:            "all namespaces with non-zero query namespace allowed",
			scope:           lifecycleAllNamespacesScope,
			queryNamespace:  "alpha",
			hasNamespace:    true,
			wantStoreCalled: true,
		},
		{
			name:            "global resource absent query namespace allowed",
			globalResource:  true,
			scope:           lifecycleAllNamespacesScope,
			wantStoreCalled: true,
		},
		{
			name:            "global resource explicit zero query namespace allowed",
			globalResource:  true,
			scope:           lifecycleAllNamespacesScope,
			queryNamespace:  "",
			hasNamespace:    true,
			wantStoreCalled: true,
		},
		{
			name:            "global resource non-zero query namespace rejected",
			globalResource:  true,
			scope:           lifecycleAllNamespacesScope,
			queryNamespace:  "alpha",
			hasNamespace:    true,
			wantReason:      ErrorReasonInvalidQueryScope,
			wantStoreCalled: false,
		},
		{
			name:            "namespaced resource absent query namespace allowed",
			scope:           lifecycleAllNamespacesScope,
			wantStoreCalled: true,
		},
		{
			name:            "namespaced resource non-zero query namespace allowed",
			scope:           lifecycleAllNamespacesScope,
			queryNamespace:  "alpha",
			hasNamespace:    true,
			wantStoreCalled: true,
		},
		{
			name:            "namespaced resource explicit zero query namespace rejected",
			scope:           lifecycleAllNamespacesScope,
			queryNamespace:  "",
			hasNamespace:    true,
			wantReason:      ErrorReasonInvalidQueryScope,
			wantStoreCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &trackingListStore{}
			options := []Option{WithStore(store)}
			if tt.globalResource {
				options = append(options, WithResourceResolver(testCatalog(t, globalTestDefinition())))
			}
			executor := testExecutor(t, options...)

			query := objectquery.Query{}
			if tt.hasNamespace {
				query.Identity.Namespace = mustLifecycleNamespaceEquals(t, tt.queryNamespace)
			}

			_, err := executor.List(context.Background(), ListRequest{
				Resource: testGVR(),
				Scope:    tt.scope(t),
				Query:    query,
			})
			if tt.wantReason != "" {
				requireLifecycleError(t, err, ErrInvalidRequest, tt.wantReason)
				if errors.Is(err, objectquery.ErrInvalidQuery) {
					t.Fatalf("errors.Is(err, objectquery.ErrInvalidQuery) = true")
				}
			} else {
				requireNoError(t, err)
			}
			if store.listCalled != tt.wantStoreCalled {
				t.Fatalf("store.List called = %v; want %v", store.listCalled, tt.wantStoreCalled)
			}
		})
	}
}

func TestListQueryErrorReasons(t *testing.T) {
	t.Run("plain invalid list request keeps invalid request reason", func(t *testing.T) {
		executor := testExecutor(t)

		_, err := executor.List(context.Background(), ListRequest{})

		requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidRequest)
		requireErrorIs(t, err, objectstore.ErrInvalidListRequest)
	})

	t.Run("query scope conflict uses query scope reason", func(t *testing.T) {
		store := &trackingListStore{}
		executor := testExecutor(t, WithStore(store))
		scope, err := objectstore.InNamespace("alpha")
		requireNoError(t, err)

		_, err = executor.List(context.Background(), ListRequest{
			Resource: testGVR(),
			Scope:    scope,
			Query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustLifecycleNamespaceEquals(t, "beta"),
				},
			},
		})

		requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidQueryScope)
		if errors.Is(err, objectquery.ErrInvalidQuery) {
			t.Fatalf("errors.Is(err, objectquery.ErrInvalidQuery) = true")
		}
		if store.listCalled {
			t.Fatalf("store.List was called")
		}
	})

	// Malformed objectquery internals are tested inside api/objectquery. From
	// objectlifecycle, selectors and requirements are intentionally opaque, so
	// public-constructible lifecycle tests cover valid query values that
	// contradict resource/scope constraints.
}

func TestListQueryPreservesRevision(t *testing.T) {
	items := []objectstore.ListItem{
		lifecycleQueryListItem(t, 1, "system", "api:v1", map[string]string{"env": "prod"}, nil, 1),
		lifecycleQueryListItem(t, 2, "system", "api:v2", map[string]string{"env": "qa"}, nil, 2),
		lifecycleQueryListItem(t, 3, "system", "api:v3", map[string]string{"env": "prod"}, nil, 3),
	}

	tests := []struct {
		name      string
		query     objectquery.Query
		wantNames []metaidentity.Name
	}{
		{
			name: "matches all",
			query: objectquery.Query{
				Labels: mustLifecycleLabelSelector(t, mustLifecycleLabelExists(t, "env")),
			},
			wantNames: []metaidentity.Name{"worker-1", "worker-2", "worker-3"},
		},
		{
			name: "matches some",
			query: objectquery.Query{
				Labels: mustLifecycleLabelSelector(t, mustLifecycleLabelEquals(t, "env", "prod")),
			},
			wantNames: []metaidentity.Name{"worker-1", "worker-3"},
		},
		{
			name: "matches none",
			query: objectquery.Query{
				Labels: mustLifecycleLabelSelector(t, mustLifecycleLabelEquals(t, "env", "dev")),
			},
			wantNames: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &trackingListStore{
				result: objectstore.ListResult{
					Items:    items,
					Revision: 42,
				},
			}
			executor := testExecutor(t, WithStore(store))

			result, err := executor.List(context.Background(), ListRequest{
				Resource: testGVR(),
				Scope:    objectstore.AllNamespaces(),
				Query:    tt.query,
			})
			requireNoError(t, err)

			requireLifecycleListNames(t, result, tt.wantNames...)
			if result.Revision != 42 {
				t.Fatalf("revision = %v; want 42", result.Revision)
			}
			if result.Revision.IsZero() {
				t.Fatal("revision is zero")
			}
		})
	}
}

func TestListQueryFilteringPreservesOrder(t *testing.T) {
	store := &trackingListStore{
		result: objectstore.ListResult{
			Items: []objectstore.ListItem{
				lifecycleQueryListItem(t, 1, "system", "api:v1", map[string]string{"env": "prod"}, nil, 1),
				lifecycleQueryListItem(t, 2, "system", "api:v2", map[string]string{"env": "qa"}, nil, 2),
				lifecycleQueryListItem(t, 3, "system", "api:v3", map[string]string{"env": "prod"}, nil, 3),
				lifecycleQueryListItem(t, 4, "system", "api:v4", map[string]string{"env": "prod"}, nil, 4),
			},
			Revision: 5,
		},
	}
	executor := testExecutor(t, WithStore(store))

	result, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
		Query: objectquery.Query{
			Labels: mustLifecycleLabelSelector(t, mustLifecycleLabelEquals(t, "env", "prod")),
		},
	})
	requireNoError(t, err)

	requireLifecycleListNames(t, result, "worker-1", "worker-3", "worker-4")
}

func TestListQueryNamespaceScopeConsistency(t *testing.T) {
	t.Run("namespace scope conflict", func(t *testing.T) {
		store := &trackingListStore{}
		executor := testExecutor(t, WithStore(store))
		scope, err := objectstore.InNamespace("alpha")
		requireNoError(t, err)

		_, err = executor.List(context.Background(), ListRequest{
			Resource: testGVR(),
			Scope:    scope,
			Query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustLifecycleNamespaceEquals(t, "beta"),
				},
			},
		})

		requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidQueryScope)
		if store.listCalled {
			t.Fatalf("store.List was called")
		}
	})

	t.Run("namespace scope same namespace", func(t *testing.T) {
		store := &trackingListStore{}
		executor := testExecutor(t, WithStore(store))
		scope, err := objectstore.InNamespace("alpha")
		requireNoError(t, err)

		_, err = executor.List(context.Background(), ListRequest{
			Resource: testGVR(),
			Scope:    scope,
			Query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustLifecycleNamespaceEquals(t, "alpha"),
				},
			},
		})
		requireNoError(t, err)

		if !store.listCalled {
			t.Fatalf("store.List was not called")
		}
	})

	t.Run("all namespaces with namespace query", func(t *testing.T) {
		executor := testExecutor(t)
		createObjectWithMetadata(t, executor, 1, "alpha", "api:alpha", nil, nil)
		createObjectWithMetadata(t, executor, 2, "beta", "api:beta", nil, nil)

		result, err := executor.List(context.Background(), ListRequest{
			Resource: testGVR(),
			Scope:    objectstore.AllNamespaces(),
			Query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustLifecycleNamespaceEquals(t, "alpha"),
				},
			},
		})
		requireNoError(t, err)

		requireLifecycleListNames(t, result, "worker-1")
	})
}

func TestListGlobalResourceQueryNamespaceConsistency(t *testing.T) {
	t.Run("absent namespace allowed", func(t *testing.T) {
		store := &trackingListStore{}
		executor := testExecutor(
			t,
			WithStore(store),
			WithResourceResolver(testCatalog(t, globalTestDefinition())),
		)

		_, err := executor.List(context.Background(), ListRequest{
			Resource: testGVR(),
			Scope:    objectstore.AllNamespaces(),
		})
		requireNoError(t, err)

		if !store.listCalled {
			t.Fatalf("store.List was not called")
		}
	})

	t.Run("explicit zero namespace allowed", func(t *testing.T) {
		store := &trackingListStore{}
		executor := testExecutor(
			t,
			WithStore(store),
			WithResourceResolver(testCatalog(t, globalTestDefinition())),
		)

		_, err := executor.List(context.Background(), ListRequest{
			Resource: testGVR(),
			Scope:    objectstore.AllNamespaces(),
			Query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustLifecycleNamespaceEquals(t, ""),
				},
			},
		})
		requireNoError(t, err)

		if !store.listCalled {
			t.Fatalf("store.List was not called")
		}
	})

	t.Run("non-zero namespace rejected", func(t *testing.T) {
		store := &trackingListStore{}
		executor := testExecutor(
			t,
			WithStore(store),
			WithResourceResolver(testCatalog(t, globalTestDefinition())),
		)

		_, err := executor.List(context.Background(), ListRequest{
			Resource: testGVR(),
			Scope:    objectstore.AllNamespaces(),
			Query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustLifecycleNamespaceEquals(t, "system"),
				},
			},
		})

		requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidQueryScope)
		if store.listCalled {
			t.Fatalf("store.List was called")
		}
	})
}

func TestListNamespacedResourceQueryNamespaceConsistency(t *testing.T) {
	t.Run("absent namespace allowed", func(t *testing.T) {
		store := &trackingListStore{}
		executor := testExecutor(t, WithStore(store))

		_, err := executor.List(context.Background(), ListRequest{
			Resource: testGVR(),
			Scope:    objectstore.AllNamespaces(),
		})
		requireNoError(t, err)

		if !store.listCalled {
			t.Fatalf("store.List was not called")
		}
	})

	t.Run("non-zero namespace allowed", func(t *testing.T) {
		store := &trackingListStore{}
		executor := testExecutor(t, WithStore(store))

		_, err := executor.List(context.Background(), ListRequest{
			Resource: testGVR(),
			Scope:    objectstore.AllNamespaces(),
			Query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustLifecycleNamespaceEquals(t, "system"),
				},
			},
		})
		requireNoError(t, err)

		if !store.listCalled {
			t.Fatalf("store.List was not called")
		}
	})

	t.Run("explicit zero namespace rejected", func(t *testing.T) {
		store := &trackingListStore{}
		executor := testExecutor(t, WithStore(store))

		_, err := executor.List(context.Background(), ListRequest{
			Resource: testGVR(),
			Scope:    objectstore.AllNamespaces(),
			Query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustLifecycleNamespaceEquals(t, ""),
				},
			},
		})

		requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidQueryScope)
		if store.listCalled {
			t.Fatalf("store.List was called")
		}
	})
}

func TestListQueryPreservesRevisionWhenFilteringAllItems(t *testing.T) {
	store := &trackingListStore{
		result: objectstore.ListResult{
			Items: []objectstore.ListItem{{
				Key: objectstore.MustKey(testGVR(), testName(1)),
				State: objectstore.State{
					Object:    testObjectWithDesired(1, value.StringValue("stored")),
					Ownership: objectownership.EmptyState(),
					Revision:  41,
				},
			}},
			Revision: 42,
		},
	}
	executor := testExecutor(t, WithStore(store))

	result, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
		Query: objectquery.Query{
			Labels: mustLifecycleLabelSelector(t, mustLifecycleLabelEquals(t, "env", "prod")),
		},
	})
	requireNoError(t, err)

	if result.Len() != 0 {
		t.Fatalf("len = %d; want 0", result.Len())
	}
	if result.Revision != 42 {
		t.Fatalf("revision = %v; want 42", result.Revision)
	}
}

func TestListQueryPreservesDetachedResults(t *testing.T) {
	executor := testExecutor(t)
	createObjectWithMetadata(
		t,
		executor,
		1,
		"system",
		"api:v1",
		map[string]string{"env": "prod"},
		nil,
	)

	result, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
		Query: objectquery.Query{
			Labels: mustLifecycleLabelSelector(t, mustLifecycleLabelEquals(t, "env", "prod")),
		},
	})
	requireNoError(t, err)

	if result.Len() != 1 {
		t.Fatalf("len = %d; want 1", result.Len())
	}
	result.Items[0].State.Object.ObjectMeta.Labels[labels.Key("env")] = labels.Value("mutated")
	result.Items[0].State.Object.Desired = value.StringValue("mutated")

	stored, err := executor.Get(context.Background(), GetRequest{
		Resource: testGVR(),
		Object:   testName(1),
	})
	requireNoError(t, err)

	got, ok := stored.State.Object.ObjectMeta.Labels.Get("env")
	if !ok || got != "prod" {
		t.Fatalf("stored label env = %q, %v; want prod, true", got, ok)
	}
	requireImage(t, stored.State, "api:v1")
}

func createObjectWithMetadata(
	t *testing.T,
	executor *Executor,
	index int,
	namespace metaidentity.Namespace,
	image string,
	rawLabels map[string]string,
	rawAnnotations map[string]string,
) Result {
	t.Helper()

	obj := testObject(index, image)
	obj.ObjectMeta.Namespace = namespace
	obj.ObjectMeta.Labels = mustLifecycleLabels(t, rawLabels)
	obj.ObjectMeta.Annotations = mustLifecycleAnnotations(t, rawAnnotations)

	result, err := executor.Create(
		context.Background(),
		CreateRequest{Object: obj, Owner: owner("creator")},
	)
	requireNoError(t, err)

	return result
}

func lifecycleAllNamespacesScope(t *testing.T) objectstore.ListScope {
	t.Helper()

	return objectstore.AllNamespaces()
}

func lifecycleNamespaceScope(namespace metaidentity.Namespace) func(t *testing.T) objectstore.ListScope {
	return func(t *testing.T) objectstore.ListScope {
		t.Helper()

		scope, err := objectstore.InNamespace(namespace)
		requireNoError(t, err)

		return scope
	}
}

func lifecycleQueryListItem(
	t *testing.T,
	index int,
	namespace metaidentity.Namespace,
	image string,
	rawLabels map[string]string,
	rawAnnotations map[string]string,
	revision objectstore.Revision,
) objectstore.ListItem {
	t.Helper()

	obj := testObjectWithDesired(index, value.StringValue(image))
	obj.ObjectMeta.Namespace = namespace
	obj.ObjectMeta.Labels = mustLifecycleLabels(t, rawLabels)
	obj.ObjectMeta.Annotations = mustLifecycleAnnotations(t, rawAnnotations)

	return objectstore.ListItem{
		Key: objectstore.MustKey(testGVR(), obj.ObjectMeta.ObjectName()),
		State: objectstore.State{
			Object:    obj,
			Ownership: objectownership.EmptyState(),
			Revision:  revision,
		},
	}
}

func mustLifecycleLabels(t *testing.T, values map[string]string) labels.Set {
	t.Helper()

	set, err := labels.FromStrings(values)
	requireNoError(t, err)

	return set
}

func mustLifecycleAnnotations(t *testing.T, values map[string]string) annotations.Set {
	t.Helper()

	set, err := annotations.FromStrings(values)
	requireNoError(t, err)

	return set
}

func mustLifecycleLabelEquals(t *testing.T, key string, val string) objectquery.LabelRequirement {
	t.Helper()

	requirement, err := objectquery.LabelEquals(key, val)
	requireNoError(t, err)

	return requirement
}

func mustLifecycleLabelExists(t *testing.T, key string) objectquery.LabelRequirement {
	t.Helper()

	requirement, err := objectquery.LabelExists(key)
	requireNoError(t, err)

	return requirement
}

func mustLifecycleLabelSelector(
	t *testing.T,
	requirements ...objectquery.LabelRequirement,
) objectquery.LabelSelector {
	t.Helper()

	selector, err := objectquery.NewLabelSelector(requirements...)
	requireNoError(t, err)

	return selector
}

func mustLifecycleAnnotationEquals(t *testing.T, key string, val string) objectquery.AnnotationRequirement {
	t.Helper()

	requirement, err := objectquery.AnnotationEquals(key, val)
	requireNoError(t, err)

	return requirement
}

func mustLifecycleAnnotationSelector(
	t *testing.T,
	requirements ...objectquery.AnnotationRequirement,
) objectquery.AnnotationSelector {
	t.Helper()

	selector, err := objectquery.NewAnnotationSelector(requirements...)
	requireNoError(t, err)

	return selector
}

func mustLifecycleNamespaceEquals(
	t *testing.T,
	namespace metaidentity.Namespace,
) objectquery.NamespaceRequirement {
	t.Helper()

	requirement, err := objectquery.NamespaceEquals(namespace)
	requireNoError(t, err)

	return requirement
}

func mustLifecycleNameEquals(t *testing.T, name metaidentity.Name) objectquery.NameRequirement {
	t.Helper()

	requirement, err := objectquery.NameEquals(name)
	requireNoError(t, err)

	return requirement
}

func requireLifecycleListNames(t *testing.T, result ListResult, names ...metaidentity.Name) {
	t.Helper()

	if result.Len() != len(names) {
		t.Fatalf("len = %d; want %d", result.Len(), len(names))
	}
	for i, name := range names {
		if result.Items[i].Key.Object.Name != name {
			t.Fatalf("item[%d].name = %s; want %s", i, result.Items[i].Key.Object.Name, name)
		}
	}
}
