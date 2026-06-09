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
	"testing"

	"arcoris.dev/apimachinery/api/objectapply"
	"arcoris.dev/apimachinery/api/objectownership"
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
	requireImage(t, result.State, "api:v1")
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
	requireImage(t, result.State, "api:v2")
	if result.Apply.Object.Desired.IsZero() {
		t.Fatalf("Apply metadata is empty")
	}
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
	requireOwnedPath(t, result.State.Ownership, owner("creator"), objectownership.Path("$.image"))
	requireOwnedPath(t, result.State.Ownership, owner("observer"), objectownership.Path("$.image"))
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
	_, err := executor.Create(
		context.Background(),
		CreateRequest{Object: testObservedObject(1, "api:v1", "true"), Owner: owner("creator")},
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
	if result.State.Object.ObjectMeta.ResourceVersion != "" {
		t.Fatalf("ResourceVersion = %q; want empty because lifecycle does not stamp metadata", result.State.Object.ObjectMeta.ResourceVersion)
	}
}

func resourceObserved() resource.VersionOption {
	return resource.Observed(observedDescriptor())
}
