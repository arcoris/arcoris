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

	"arcoris.dev/apimachinery/api/objectapply"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/value"
)

func TestApplyMissingObjectCreatesState(t *testing.T) {
	executor := testExecutor(t)

	result, err := executor.Apply(
		context.Background(),
		ApplyRequest{Object: testObject(1, "api:v1"), Owner: owner("creator")},
	)
	requireNoError(t, err)

	requireEffect(t, result.Result, OperationApply, EffectCreated)
	if result.Revision != result.State.Revision {
		t.Fatalf("result revision = %v; want state revision %v", result.Revision, result.State.Revision)
	}
	requireImage(t, result.State, "api:v1")
	requireNormalizedOwnership(t, result.State.Ownership)
}

func TestApplyExistingObjectCommitsUpdate(t *testing.T) {
	executor := testExecutor(t)
	created := createObject(t, executor, 1, "api:v1", owner("creator"))

	result, err := executor.Apply(
		context.Background(),
		ApplyRequest{Object: testObject(1, "api:v2"), Owner: owner("creator")},
	)
	requireNoError(t, err)

	requireEffect(t, result.Result, OperationApply, EffectUpdated)
	if !created.State.Revision.Before(result.State.Revision) {
		t.Fatalf("revision = %v; want after %v", result.State.Revision, created.State.Revision)
	}
	if result.Revision != result.State.Revision {
		t.Fatalf("result revision = %v; want state revision %v", result.Revision, result.State.Revision)
	}
	requireImage(t, result.State, "api:v2")
	if result.Apply.Object.Desired.IsZero() {
		t.Fatalf("Apply metadata is empty")
	}
	requireNormalizedOwnership(t, result.State.Ownership)
}

func TestApplySameObjectCurrentlyCommitsUpdate(t *testing.T) {
	executor := testExecutor(t)
	created := createObject(t, executor, 1, "api:v1", owner("creator"))

	result, err := executor.Apply(
		context.Background(),
		ApplyRequest{Object: testObject(1, "api:v1"), Owner: owner("creator")},
	)
	requireNoError(t, err)

	requireEffect(t, result.Result, OperationApply, EffectUpdated)
	if !created.State.Revision.Before(result.State.Revision) {
		t.Fatalf("revision = %v; want after %v", result.State.Revision, created.State.Revision)
	}
}

func TestApplySameValueDifferentOwnerChangesOwnership(t *testing.T) {
	executor := testExecutor(t)
	createObject(t, executor, 1, "api:v1", owner("creator"))

	result, err := executor.Apply(
		context.Background(),
		ApplyRequest{Object: testObject(1, "api:v1"), Owner: owner("observer")},
	)
	requireNoError(t, err)

	requireEffect(t, result.Result, OperationApply, EffectUpdated)
	requireOwnedPath(t, result.State.Ownership, owner("creator"), ownershipField("$.image"))
	requireOwnedPath(t, result.State.Ownership, owner("observer"), ownershipField("$.image"))
}

func TestApplyInvalidObjectReturnsValidationFailed(t *testing.T) {
	executor := testExecutor(t)
	obj := testObject(1, "api:v1")
	obj.Desired = value.StringValue("not-object")

	_, err := executor.Apply(
		context.Background(),
		ApplyRequest{Object: obj, Owner: owner("creator")},
	)

	requireLifecycleError(t, err, ErrValidationFailed, ErrorReasonValidationFailed)
}

func TestApplyMissingObjectRejectsObserved(t *testing.T) {
	executor := testExecutor(
		t,
		WithResourceResolver(testCatalog(t, testDefinition(resourceObserved()))),
	)

	_, err := executor.Apply(
		context.Background(),
		ApplyRequest{Object: testObservedObject(1, "api:v1", "true"), Owner: owner("creator")},
	)

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonUnsupportedObservedApply)
	requireErrorIs(t, err, objectapply.ErrUnsupportedObservedApply)
}

func TestApplyExistingObjectRejectsObserved(t *testing.T) {
	executor := testExecutor(
		t,
		WithResourceResolver(testCatalog(t, testDefinition(resourceObserved()))),
	)
	_, err := executor.Create(
		context.Background(),
		CreateRequest{Object: testObservedObject(1, "api:v1", "true"), Owner: owner("creator")},
	)
	requireNoError(t, err)

	_, err = executor.Apply(
		context.Background(),
		ApplyRequest{Object: testObservedObject(1, "api:v2", "false"), Owner: owner("creator")},
	)

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonUnsupportedObservedApply)
	requireErrorIs(t, err, objectapply.ErrUnsupportedObservedApply)
}

func TestApplyRejectsZeroObservedPointer(t *testing.T) {
	executor := testExecutor(
		t,
		WithResourceResolver(testCatalog(t, testDefinition(resourceObserved()))),
	)
	obj := testObject(1, "api:v1")
	observed := value.Value{}
	obj.Observed = &observed

	_, err := executor.Apply(
		context.Background(),
		ApplyRequest{Object: obj, Owner: owner("creator")},
	)

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonUnsupportedObservedApply)
	requireErrorIs(t, err, objectapply.ErrUnsupportedObservedApply)
}

func TestApplyRejectsObservedBeforeStoreGet(t *testing.T) {
	store := &mustNotGetStore{}
	executor := testExecutor(
		t,
		WithStore(store),
		WithResourceResolver(testCatalog(t, testDefinition(resourceObserved()))),
	)

	_, err := executor.Apply(
		context.Background(),
		ApplyRequest{Object: testObservedObject(1, "api:v1", "true"), Owner: owner("creator")},
	)

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonUnsupportedObservedApply)
	if store.getCalled {
		t.Fatalf("store Get was called before rejecting observed apply")
	}
}

func TestApplyOwnershipConflictMapsToConflict(t *testing.T) {
	executor := testExecutor(t)
	createObject(t, executor, 1, "api:v1", owner("creator"))

	_, err := executor.Apply(
		context.Background(),
		ApplyRequest{Object: testObject(1, "api:v2"), Owner: owner("other")},
	)

	requireLifecycleError(t, err, ErrConflict, ErrorReasonConflict)
	requireErrorIs(t, err, objectapply.ErrConflict)
}

func TestApplyForceResolvesOwnershipConflict(t *testing.T) {
	executor := testExecutor(t)
	createObject(t, executor, 1, "api:v1", owner("creator"))

	result, err := executor.Apply(
		context.Background(),
		ApplyRequest{Object: testObject(1, "api:v2"), Owner: owner("other"), Force: true},
	)
	requireNoError(t, err)

	requireEffect(t, result.Result, OperationApply, EffectUpdated)
	requireImage(t, result.State, "api:v2")
}

func TestApplyPreservesLiveObservedAndMetadata(t *testing.T) {
	executor := testExecutor(
		t,
		WithResourceResolver(testCatalog(t, testDefinition(resourceObserved()))),
	)
	created, err := executor.Create(
		context.Background(),
		CreateRequest{Object: testObject(1, "api:v1"), Owner: owner("creator")},
	)
	requireNoError(t, err)

	observed, err := executor.UpdateObserved(
		context.Background(),
		UpdateObservedRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Observed: objectValue(member("ready", value.StringValue("true"))),
			Owner:    owner("controller"),
			Expected: created.State.Revision,
		},
	)
	requireNoError(t, err)

	patched, err := executor.PatchMetadata(
		context.Background(),
		PatchMetadataRequest{
			Resource:    testGVR(),
			Object:      testName(1),
			Labels:      map[string]*string{"app": stringPtr("worker")},
			Annotations: map[string]*string{"with.dots": stringPtr("yes")},
			Owner:       owner("metadata"),
			Expected:    observed.State.Revision,
		},
	)
	requireNoError(t, err)

	result, err := executor.Apply(
		context.Background(),
		ApplyRequest{Object: testObject(1, "api:v2"), Owner: owner("creator")},
	)
	requireNoError(t, err)

	if result.State.Object.Observed == nil {
		t.Fatalf("Observed missing")
	}
	requireObservedReady(t, result.State, "true")
	requireLabel(t, result.Result, "app", "worker")
	requireAnnotation(t, result.Result, "with.dots", "yes")
	if result.State.Object.ObjectMeta.ResourceVersion != "" {
		t.Fatalf("ResourceVersion = %q; want empty because lifecycle does not stamp metadata", result.State.Object.ObjectMeta.ResourceVersion)
	}
	if !patched.State.Revision.Before(result.State.Revision) {
		t.Fatalf("revision = %v; want after %v", result.State.Revision, patched.State.Revision)
	}
	requireOwnedPath(t, result.State.Ownership, owner("creator"), ownershipField("$.image"))
	requireSurfaceOwnedPath(t, result.State.Ownership.Observed(), owner("controller"), ownershipField("$.ready"))
	requireSurfaceOwnedPath(t, result.State.Ownership.Metadata().Labels(), owner("metadata"), ownershipField(`$["app"]`))
	requireSurfaceOwnedPath(t, result.State.Ownership.Metadata().Annotations(), owner("metadata"), ownershipField(`$["with.dots"]`))
}

func resourceObserved() resource.VersionOption {
	return resource.Observed(observedDescriptor())
}

type mustNotGetStore struct {
	getCalled bool
}

func (s *mustNotGetStore) Get(context.Context, objectstore.Key) (objectstore.State, bool, error) {
	s.getCalled = true

	return objectstore.State{}, false, errors.New("unexpected get")
}

func (s *mustNotGetStore) Create(context.Context, objectstore.Key, objectstore.State) (objectstore.State, error) {
	return objectstore.State{}, errors.New("unexpected create")
}

func (s *mustNotGetStore) Update(context.Context, objectstore.Key, objectstore.Revision, objectstore.State) (objectstore.State, error) {
	return objectstore.State{}, errors.New("unexpected update")
}

func (s *mustNotGetStore) Delete(context.Context, objectstore.Key, objectstore.Revision) (objectstore.DeleteResult, error) {
	return objectstore.DeleteResult{}, errors.New("unexpected delete")
}
