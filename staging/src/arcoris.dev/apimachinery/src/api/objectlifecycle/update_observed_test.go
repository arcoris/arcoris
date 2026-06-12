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

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestUpdateObservedReplacesObservedOnly(t *testing.T) {
	executor := testExecutor(
		t,
		WithResourceResolver(testCatalog(t, testDefinition(resourceObserved()))),
	)
	created := createObject(t, executor, 1, "api:v1", owner("creator"))

	result, err := executor.UpdateObserved(
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

	requireEffect(t, result, OperationUpdateObserved, EffectUpdated)
	if !created.State.Revision.Before(result.State.Revision) {
		t.Fatalf("revision = %v; want after %v", result.State.Revision, created.State.Revision)
	}
	if result.Revision != result.State.Revision {
		t.Fatalf("result revision = %v; want state revision %v", result.Revision, result.State.Revision)
	}
	requireImage(t, result.State, "api:v1")
	requireObservedReady(t, result.State, "true")
	requireOwnedPath(t, result.State.Ownership, owner("creator"), ownershipField("$.image"))
	requireSurfaceOwnedPath(
		t,
		result.State.Ownership.Observed(),
		owner("controller"),
		ownershipPath(fieldpath.Root().Field(fieldpath.MustFieldName("ready"))),
	)
	if !result.State.Ownership.Metadata().IsEmpty() {
		t.Fatalf("metadata ownership = %#v; want empty", result.State.Ownership.Metadata())
	}
}

func TestUpdateObservedRequiresObservedDescriptor(t *testing.T) {
	executor := testExecutor(t)
	created := createObject(t, executor, 1, "api:v1", owner("creator"))

	_, err := executor.UpdateObserved(
		context.Background(),
		UpdateObservedRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Observed: objectValue(member("ready", value.StringValue("true"))),
			Owner:    owner("controller"),
			Expected: created.State.Revision,
		},
	)

	requireLifecycleError(t, err, ErrValidationFailed, ErrorReasonObservedNotDefined)
}

func TestUpdateObservedRejectsInvalidObservedPayload(t *testing.T) {
	executor := testExecutor(
		t,
		WithResourceResolver(testCatalog(t, testDefinition(resourceObserved()))),
	)
	created := createObject(t, executor, 1, "api:v1", owner("creator"))

	_, err := executor.UpdateObserved(
		context.Background(),
		UpdateObservedRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Observed: value.StringValue("not-object"),
			Owner:    owner("controller"),
			Expected: created.State.Revision,
		},
	)

	requireLifecycleError(t, err, ErrValidationFailed, ErrorReasonInvalidObserved)
}

func TestUpdateObservedMissingObjectReturnsNotFound(t *testing.T) {
	executor := testExecutor(
		t,
		WithResourceResolver(testCatalog(t, testDefinition(resourceObserved()))),
	)

	_, err := executor.UpdateObserved(
		context.Background(),
		UpdateObservedRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Observed: objectValue(member("ready", value.StringValue("true"))),
			Owner:    owner("controller"),
			Expected: 1,
		},
	)

	requireLifecycleError(t, err, ErrNotFound, ErrorReasonNotFound)
}

func TestUpdateObservedStaleRevisionReturnsStaleRevision(t *testing.T) {
	executor := testExecutor(
		t,
		WithResourceResolver(testCatalog(t, testDefinition(resourceObserved()))),
	)
	createObject(t, executor, 1, "api:v1", owner("creator"))

	_, err := executor.UpdateObserved(
		context.Background(),
		UpdateObservedRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Observed: objectValue(member("ready", value.StringValue("true"))),
			Owner:    owner("controller"),
			Expected: 99,
		},
	)

	requireLifecycleError(t, err, ErrStaleRevision, ErrorReasonStaleRevision)
}

func TestUpdateObservedRequiresObservedValidatorOnlyWhenUsed(t *testing.T) {
	executor, err := NewExecutor(
		WithStore(testStore(t)),
		WithResourceResolver(testCatalog(t, testDefinition(resourceObserved()))),
		WithDesiredValidator(valuevalidation.SurfaceValidator{}),
	)
	requireNoError(t, err)

	_, err = executor.UpdateObserved(
		context.Background(),
		UpdateObservedRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Observed: objectValue(member("ready", value.StringValue("true"))),
			Owner:    owner("controller"),
			Expected: 1,
		},
	)

	requireLifecycleError(t, err, ErrInvalidExecutor, ErrorReasonInvalidExecutor)
	requireErrorIs(t, err, ErrNilObservedValidator)
}

func TestUpdateObservedMissingResourceReturnsResourceNotFound(t *testing.T) {
	executor := testExecutor(
		t,
		WithResourceResolver(testCatalog(t, testDefinition(resourceObserved()))),
	)
	req := UpdateObservedRequest{
		Resource: testGVR(),
		Object:   testName(1),
		Observed: objectValue(member("ready", value.StringValue("true"))),
		Owner:    owner("controller"),
		Expected: 1,
	}
	req.Resource.Resource = "unknowns"

	_, err := executor.UpdateObserved(context.Background(), req)

	requireLifecycleError(t, err, ErrResourceNotFound, ErrorReasonResourceNotFound)
}
