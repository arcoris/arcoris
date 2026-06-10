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
	"arcoris.dev/apimachinery/api/meta/stamp"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestCreateValidObjectStoresCommittedState(t *testing.T) {
	executor := testExecutor(t)

	result, err := executor.Create(
		context.Background(),
		CreateRequest{Object: testObject(1, "api:v1"), Owner: owner("creator")},
	)
	requireNoError(t, err)

	requireEffect(t, result, OperationCreate, EffectCreated)
	if !result.State.Revision.IsValid() {
		t.Fatalf("revision is invalid")
	}
	requireImage(t, result.State, "api:v1")
	requireOwnedPath(t, result.State.Ownership, owner("creator"), objectownership.Path("$.image"))
}

func TestCreateExistingReturnsAlreadyExists(t *testing.T) {
	executor := testExecutor(t)
	createObject(t, executor, 1, "api:v1", owner("creator"))

	_, err := executor.Create(
		context.Background(),
		CreateRequest{Object: testObject(1, "api:v1"), Owner: owner("creator")},
	)

	requireLifecycleError(t, err, ErrAlreadyExists, ErrorReasonAlreadyExists)
}

func TestCreateInvalidObjectReturnsValidationFailed(t *testing.T) {
	executor := testExecutor(t)
	obj := testObject(1, "api:v1")
	obj.Desired = value.StringValue("not-object")

	_, err := executor.Create(
		context.Background(),
		CreateRequest{Object: obj, Owner: owner("creator")},
	)

	requireLifecycleError(t, err, ErrValidationFailed, ErrorReasonValidationFailed)
}

func TestCreateMissingResourceReturnsResourceNotFound(t *testing.T) {
	executor := testExecutor(t, WithResourceResolver(testCatalog(t, testDefinition())))
	obj := testObject(1, "api:v1")
	obj.TypeMeta.Kind = "Unknown"

	_, err := executor.Create(
		context.Background(),
		CreateRequest{Object: obj, Owner: owner("creator")},
	)

	requireLifecycleError(t, err, ErrResourceNotFound, ErrorReasonResourceNotFound)
}

func TestCreateDoesNotStampMetadata(t *testing.T) {
	executor := testExecutor(t)

	result, err := executor.Create(
		context.Background(),
		CreateRequest{Object: testObject(1, "api:v1"), Owner: owner("creator")},
	)
	requireNoError(t, err)

	if result.State.Object.ObjectMeta.ResourceVersion != "" {
		t.Fatalf("ResourceVersion = %q; want empty", result.State.Object.ObjectMeta.ResourceVersion)
	}
	if result.State.Object.ObjectMeta.Generation != stamp.Generation(0) {
		t.Fatalf("Generation = %d; want 0", result.State.Object.ObjectMeta.Generation)
	}
	if result.State.Object.ObjectMeta.UID != "" {
		t.Fatalf("UID = %q; want empty", result.State.Object.ObjectMeta.UID)
	}
}

func TestCreateInitializesNestedDesiredOwnership(t *testing.T) {
	executor := testExecutor(t, WithResourceResolver(testCatalog(t, ownershipDefinition())))
	obj := testObjectWithDesired(1, objectValue(
		member("template", objectValue(member("image", value.StringValue("api:v1")))),
	))

	result, err := executor.Create(
		context.Background(),
		CreateRequest{Object: obj, Owner: owner("creator")},
	)
	requireNoError(t, err)

	requireOwnedPath(
		t,
		result.State.Ownership,
		owner("creator"),
		ownershipPath(fieldpath.Root().Field(fieldpath.MustFieldName("template")).Field(fieldpath.MustFieldName("image"))),
	)
}

func TestCreateInitializesListMapDesiredOwnership(t *testing.T) {
	executor := testExecutor(t, WithResourceResolver(testCatalog(t, ownershipDefinition())))
	selector := fieldpath.MustSelector(fieldpath.NewSelectorEntry("type", fieldpath.StringLiteral("Ready")))
	obj := testObjectWithDesired(1, objectValue(
		member("conditions", value.MustListValue(objectValue(
			member("type", value.StringValue("Ready")),
			member("detail", objectValue(member("message", value.StringValue("ok")))),
		))),
	))

	result, err := executor.Create(
		context.Background(),
		CreateRequest{Object: obj, Owner: owner("creator")},
	)
	requireNoError(t, err)

	requireOwnedPath(
		t,
		result.State.Ownership,
		owner("creator"),
		ownershipPath(fieldpath.Root().Field(fieldpath.MustFieldName("conditions")).Select(selector).Field(fieldpath.MustFieldName("detail")).Field(fieldpath.MustFieldName("message"))),
	)
}

func ownershipDefinition() resource.Definition {
	condition := types.Object(
		types.Field("type").String().Required(),
		types.Field("detail").Object(
			types.Field("message").String().Required(),
		).Required(),
	)
	desired := types.Object(
		types.Field("template").Object(
			types.Field("image").String().Required(),
		).Optional(),
		types.Field("conditions").ListOf(condition).Optional().Map("type"),
	).Descriptor()

	return resource.NewDefinition(
		testGroup,
		"Worker",
		"workers",
		resource.ScopeNamespaced,
		resource.NewVersion("v1", desired, resource.Exposed(), resource.Canonical()),
	)
}
