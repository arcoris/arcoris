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

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/value"
)

func TestListReturnsCommittedItems(t *testing.T) {
	executor := testExecutor(t)
	first := createObject(t, executor, 1, "api:v1", owner("creator"))
	second := createObject(t, executor, 2, "api:v2", owner("creator"))

	result, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
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

func TestListNamespaceScope(t *testing.T) {
	executor := testExecutor(t)
	createObject(t, executor, 1, "api:v1", owner("creator"))
	scope, err := objectstore.InNamespace("other")
	requireNoError(t, err)

	result, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    scope,
	})
	requireNoError(t, err)

	if result.Len() != 0 {
		t.Fatalf("len = %d; want 0", result.Len())
	}
}

func TestListGlobalResourceAllowsAllNamespacesScope(t *testing.T) {
	executor := testExecutor(t, WithResourceResolver(testCatalog(t, globalTestDefinition())))

	result, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
	})
	requireNoError(t, err)

	if result.Len() != 0 {
		t.Fatalf("len = %d; want 0", result.Len())
	}
}

func TestListUnknownResourceReturnsResourceNotFound(t *testing.T) {
	executor := testExecutor(t)
	req := ListRequest{Resource: testGVR(), Scope: objectstore.AllNamespaces()}
	req.Resource.Resource = "unknowns"

	_, err := executor.List(context.Background(), req)

	requireLifecycleError(t, err, ErrResourceNotFound, ErrorReasonResourceNotFound)
}

func TestListUnknownResourceDoesNotCallStore(t *testing.T) {
	store := &trackingListStore{}
	executor := testExecutor(t, WithStore(store))
	req := ListRequest{Resource: testGVR(), Scope: objectstore.AllNamespaces()}
	req.Resource.Resource = "unknowns"

	_, err := executor.List(context.Background(), req)

	requireLifecycleError(t, err, ErrResourceNotFound, ErrorReasonResourceNotFound)
	if store.listCalled {
		t.Fatalf("store.List was called")
	}
}

func TestListRejectsInvalidRequest(t *testing.T) {
	executor := testExecutor(t)

	_, err := executor.List(context.Background(), ListRequest{})

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidRequest)
	requireErrorIs(t, err, objectstore.ErrInvalidListRequest)
}

func TestListRejectsInvalidRequestBeforeStoreCall(t *testing.T) {
	store := &trackingListStore{}
	executor := testExecutor(t, WithStore(store))

	_, err := executor.List(context.Background(), ListRequest{})

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidRequest)
	if store.listCalled {
		t.Fatalf("store.List was called")
	}
}

func TestListRejectsNamespaceScopeForGlobalResource(t *testing.T) {
	executor := testExecutor(t, WithResourceResolver(testCatalog(t, globalTestDefinition())))
	scope, err := objectstore.InNamespace("system")
	requireNoError(t, err)

	_, err = executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    scope,
	})

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidRequest)
	requireErrorIs(t, err, objectstore.ErrInvalidListRequest)
}

func TestListRejectsInvalidResourceContract(t *testing.T) {
	executor := testExecutor(t, WithResourceResolver(inconsistentListResourceResolver{}))

	_, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
	})

	requireLifecycleError(t, err, ErrValidationFailed, ErrorReasonInvalidResourceContract)
	requireErrorIs(t, err, ErrInvalidResourceContract)
}

func TestListPassesResolvedGVRToStore(t *testing.T) {
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
	if store.request.Resource != testGVR() {
		t.Fatalf("store resource = %s; want %s", store.request.Resource, testGVR())
	}
}

func TestListPropagatesStoreErrors(t *testing.T) {
	executor := testExecutor(t, WithStore(listErrorStore{err: objectstore.ErrInvalidListRequest}))

	_, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
	})

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidRequest)
	requireErrorIs(t, err, objectstore.ErrInvalidListRequest)
}

func TestListReturnsDetachedStoreResult(t *testing.T) {
	key := objectstore.MustKey(testGVR(), testName(1))
	store := &trackingListStore{
		result: objectstore.ListResult{
			Items: []objectstore.ListItem{{
				Key: key,
				State: objectstore.State{
					Object:    testObjectWithDesired(1, value.StringValue("stored")),
					Ownership: objectownership.EmptyState(),
					Revision:  1,
				},
			}},
			Revision: 1,
		},
	}
	executor := testExecutor(t, WithStore(store))

	first, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
	})
	requireNoError(t, err)
	first.Items[0].State.Object.Desired = value.StringValue("mutated")
	first.Items[0] = objectstore.ListItem{}

	second, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
	})
	requireNoError(t, err)

	if second.Len() != 1 {
		t.Fatalf("len = %d; want 1", second.Len())
	}
	got, ok := second.Items[0].State.Object.Desired.AsString()
	if !ok || got != "stored" {
		t.Fatalf("desired = %q, %v; want stored, true", got, ok)
	}
}

func TestListDoesNotValidateStoredDesiredPayload(t *testing.T) {
	store := testStore(t)
	key := objectstore.MustKey(testGVR(), testName(1))
	created, err := store.Create(
		context.Background(),
		key,
		objectstore.State{
			Object:    testObjectWithDesired(1, value.StringValue("descriptor-invalid")),
			Ownership: objectownership.EmptyState(),
		},
	)
	requireNoError(t, err)
	executor := testExecutor(t, WithStore(store))

	result, err := executor.List(context.Background(), ListRequest{
		Resource: testGVR(),
		Scope:    objectstore.AllNamespaces(),
	})
	requireNoError(t, err)

	if result.Len() != 1 {
		t.Fatalf("len = %d; want 1", result.Len())
	}
	if result.Revision != created.Revision {
		t.Fatalf("revision = %v; want %v", result.Revision, created.Revision)
	}
	got, ok := result.Items[0].State.Object.Desired.AsString()
	if !ok || got != "descriptor-invalid" {
		t.Fatalf("desired = %q, %v; want descriptor-invalid, true", got, ok)
	}
}

type inconsistentListResourceResolver struct{}

func (inconsistentListResourceResolver) ResolveResource(
	apiidentity.GroupResource,
) (resource.Definition, bool) {
	return resource.Definition{}, false
}

func (inconsistentListResourceResolver) ResolveKind(
	apiidentity.GroupKind,
) (resource.Definition, bool) {
	return resource.Definition{}, false
}

func (inconsistentListResourceResolver) ResolveVersionResource(
	apiidentity.GroupVersionResource,
) (resource.Definition, resource.VersionDefinition, bool) {
	definition := resource.NewDefinition(
		testGroup,
		"Worker",
		"other-workers",
		resource.ScopeNamespaced,
		resource.NewVersion("v1", desiredDescriptor(), resource.Exposed(), resource.Canonical()),
	)

	return definition, resource.NewVersion("v1", desiredDescriptor()), true
}

func (inconsistentListResourceResolver) ResolveVersionKind(
	apiidentity.GroupVersionKind,
) (resource.Definition, resource.VersionDefinition, bool) {
	return resource.Definition{}, resource.VersionDefinition{}, false
}

// trackingListStore records List calls and returns a caller-owned result.
type trackingListStore struct {
	// listCalled records whether List reached the store boundary.
	listCalled bool

	// request is the last request passed to List.
	request objectstore.ListRequest

	// result is returned by List without pre-cloning so lifecycle detachment is
	// tested at the lifecycle boundary.
	result objectstore.ListResult

	// err is returned by List when set.
	err error
}

// Get is unexpected for trackingListStore because these tests only exercise List.
func (s *trackingListStore) Get(context.Context, objectstore.Key) (objectstore.State, bool, error) {
	return objectstore.State{}, false, errors.New("unexpected get")
}

// List records the request and returns the configured result or error.
func (s *trackingListStore) List(
	_ context.Context,
	request objectstore.ListRequest,
) (objectstore.ListResult, error) {
	s.listCalled = true
	s.request = request
	return s.result, s.err
}

// Create is unexpected for trackingListStore because these tests only exercise List.
func (s *trackingListStore) Create(
	context.Context,
	objectstore.Key,
	objectstore.State,
) (objectstore.State, error) {
	return objectstore.State{}, errors.New("unexpected create")
}

// Update is unexpected for trackingListStore because these tests only exercise List.
func (s *trackingListStore) Update(
	context.Context,
	objectstore.Key,
	objectstore.Revision,
	objectstore.State,
) (objectstore.State, error) {
	return objectstore.State{}, errors.New("unexpected update")
}

// Delete is unexpected for trackingListStore because these tests only exercise List.
func (s *trackingListStore) Delete(
	context.Context,
	objectstore.Key,
	objectstore.Revision,
) (objectstore.DeleteResult, error) {
	return objectstore.DeleteResult{}, errors.New("unexpected delete")
}

// globalTestDefinition returns a global resource family for scope rejection tests.
func globalTestDefinition() resource.Definition {
	return resource.NewDefinition(
		testGroup,
		"Worker",
		"workers",
		resource.ScopeGlobal,
		resource.NewVersion("v1", desiredDescriptor(), resource.Exposed(), resource.Canonical()),
	)
}

// requireLifecycleListItem checks the storage key, revision, and Desired image.
func requireLifecycleListItem(
	t *testing.T,
	item objectstore.ListItem,
	expectedKey objectstore.Key,
	revision objectstore.Revision,
	image string,
) {
	t.Helper()

	if !item.Key.Equal(expectedKey) {
		t.Fatalf("key = %s; want %s", item.Key, expectedKey)
	}
	if item.State.Revision != revision {
		t.Fatalf("revision = %v; want %v", item.State.Revision, revision)
	}
	requireImage(t, item.State, image)
}

// listErrorStore is a lifecycle test store that fails only List.
type listErrorStore struct {
	// err is returned from List to exercise lifecycle store-error mapping.
	err error
}

// Get is unexpected for listErrorStore because List should not call Get.
func (s listErrorStore) Get(context.Context, objectstore.Key) (objectstore.State, bool, error) {
	return objectstore.State{}, false, errors.New("unexpected get")
}

// List returns the configured error for store-error propagation tests.
func (s listErrorStore) List(context.Context, objectstore.ListRequest) (objectstore.ListResult, error) {
	return objectstore.ListResult{}, s.err
}

// Create is unexpected for listErrorStore because List is read-only.
func (s listErrorStore) Create(context.Context, objectstore.Key, objectstore.State) (objectstore.State, error) {
	return objectstore.State{}, errors.New("unexpected create")
}

// Update is unexpected for listErrorStore because List is read-only.
func (s listErrorStore) Update(context.Context, objectstore.Key, objectstore.Revision, objectstore.State) (objectstore.State, error) {
	return objectstore.State{}, errors.New("unexpected update")
}

// Delete is unexpected for listErrorStore because List is read-only.
func (s listErrorStore) Delete(context.Context, objectstore.Key, objectstore.Revision) (objectstore.DeleteResult, error) {
	return objectstore.DeleteResult{}, errors.New("unexpected delete")
}
