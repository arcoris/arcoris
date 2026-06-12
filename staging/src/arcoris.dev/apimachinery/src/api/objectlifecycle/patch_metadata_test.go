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
	"arcoris.dev/apimachinery/api/meta/annotations"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/meta/stamp"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

func TestPatchMetadataSetsLabelsAndAnnotationsOnly(t *testing.T) {
	executor := testExecutor(
		t,
		WithResourceResolver(testCatalog(t, testDefinition(resourceObserved()))),
	)
	created := createObject(t, executor, 1, "api:v1", owner("creator"))
	observed := updateObservedForPatchTest(t, executor, created.State.Revision)

	result, err := executor.PatchMetadata(
		context.Background(),
		PatchMetadataRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Labels: map[string]*string{
				"app":                        stringPtr("worker"),
				"scheduler.arcoris.dev/mode": stringPtr("active"),
			},
			Annotations: map[string]*string{
				"with.dots":                  stringPtr("yes"),
				"scheduler.arcoris.dev/mode": stringPtr("active"),
			},
			Owner:    owner("metadata"),
			Expected: observed.State.Revision,
		},
	)
	requireNoError(t, err)

	requireEffect(t, result, OperationPatchMetadata, EffectUpdated)
	if !observed.State.Revision.Before(result.State.Revision) {
		t.Fatalf("revision = %v; want after %v", result.State.Revision, observed.State.Revision)
	}
	requireImage(t, result.State, "api:v1")
	requireObservedReady(t, result.State, "true")
	requireLabel(t, result, "app", "worker")
	requireLabel(t, result, "scheduler.arcoris.dev/mode", "active")
	requireAnnotation(t, result, "with.dots", "yes")
	requireAnnotation(t, result, "scheduler.arcoris.dev/mode", "active")
	if result.State.Object.ObjectMeta.Name != testName(1).Name ||
		result.State.Object.ObjectMeta.Namespace != testName(1).Namespace {
		t.Fatalf("object identity changed: %#v", result.State.Object.ObjectMeta.ObjectName())
	}
	if result.State.Object.ObjectMeta.ResourceVersion != "" {
		t.Fatalf("ResourceVersion = %q; want empty", result.State.Object.ObjectMeta.ResourceVersion)
	}
	if result.State.Object.ObjectMeta.Generation != stamp.Generation(0) {
		t.Fatalf("Generation = %d; want zero", result.State.Object.ObjectMeta.Generation)
	}
	requireOwnedPath(t, result.State.Ownership, owner("creator"), ownershipField("$.image"))
	requireSurfaceOwnedPath(
		t,
		result.State.Ownership.Observed(),
		owner("controller"),
		ownershipPath(fieldpath.Root().Field(fieldpath.MustFieldName("ready"))),
	)
	requireSurfaceOwnedPath(
		t,
		result.State.Ownership.Metadata().Labels(),
		owner("metadata"),
		ownershipPath(fieldpath.Root().Key(fieldpath.MustMapKey("scheduler.arcoris.dev/mode"))),
	)
	requireSurfaceOwnedPath(
		t,
		result.State.Ownership.Metadata().Annotations(),
		owner("metadata"),
		ownershipPath(fieldpath.Root().Key(fieldpath.MustMapKey("with.dots"))),
	)
}

func TestPatchMetadataDeletesKeys(t *testing.T) {
	executor := testExecutor(t)
	created := createObject(t, executor, 1, "api:v1", owner("creator"))
	setResult, err := executor.PatchMetadata(
		context.Background(),
		PatchMetadataRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Labels:   map[string]*string{"app": stringPtr("worker")},
			Owner:    owner("metadata"),
			Expected: created.State.Revision,
		},
	)
	requireNoError(t, err)

	result, err := executor.PatchMetadata(
		context.Background(),
		PatchMetadataRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Labels:   map[string]*string{"app": nil},
			Owner:    owner("cleanup"),
			Expected: setResult.State.Revision,
		},
	)
	requireNoError(t, err)

	if result.State.Object.ObjectMeta.Labels.Has(labels.Key("app")) {
		t.Fatalf("label app still present: %#v", result.State.Object.ObjectMeta.Labels)
	}
	requireSurfaceOwnedPath(
		t,
		result.State.Ownership.Metadata().Labels(),
		owner("cleanup"),
		ownershipPath(fieldpath.Root().Key(fieldpath.MustMapKey("app"))),
	)
}

func TestPatchMetadataRejectsEmptyPatch(t *testing.T) {
	executor := testExecutor(t)

	_, err := executor.PatchMetadata(
		context.Background(),
		PatchMetadataRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Owner:    owner("metadata"),
			Expected: 1,
		},
	)

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonEmptyMetadataPatch)
}

func TestPatchMetadataRejectsInvalidLabelKey(t *testing.T) {
	executor := testExecutor(t)

	_, err := executor.PatchMetadata(
		context.Background(),
		PatchMetadataRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Labels:   map[string]*string{"Role": stringPtr("worker")},
			Owner:    owner("metadata"),
			Expected: 1,
		},
	)

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidMetadataKey)
	requireErrorIs(t, err, labels.ErrInvalidKey)
}

func TestPatchMetadataRejectsInvalidLabelValue(t *testing.T) {
	executor := testExecutor(t)

	_, err := executor.PatchMetadata(
		context.Background(),
		PatchMetadataRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Labels:   map[string]*string{"app": stringPtr("bad value")},
			Owner:    owner("metadata"),
			Expected: 1,
		},
	)

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidMetadataPatch)
	requireErrorIs(t, err, labels.ErrInvalidValue)
}

func TestPatchMetadataRejectsInvalidAnnotationValue(t *testing.T) {
	executor := testExecutor(t)

	_, err := executor.PatchMetadata(
		context.Background(),
		PatchMetadataRequest{
			Resource:    testGVR(),
			Object:      testName(1),
			Annotations: map[string]*string{"note": stringPtr("bad\nvalue")},
			Owner:       owner("metadata"),
			Expected:    1,
		},
	)

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidMetadataPatch)
	requireErrorIs(t, err, annotations.ErrInvalidValue)
}

func TestPatchMetadataMissingObjectReturnsNotFound(t *testing.T) {
	executor := testExecutor(t)

	_, err := executor.PatchMetadata(
		context.Background(),
		PatchMetadataRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Labels:   map[string]*string{"app": stringPtr("worker")},
			Owner:    owner("metadata"),
			Expected: 1,
		},
	)

	requireLifecycleError(t, err, ErrNotFound, ErrorReasonNotFound)
}

func TestPatchMetadataStaleRevisionReturnsStaleRevision(t *testing.T) {
	executor := testExecutor(t)
	createObject(t, executor, 1, "api:v1", owner("creator"))

	_, err := executor.PatchMetadata(
		context.Background(),
		PatchMetadataRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Labels:   map[string]*string{"app": stringPtr("worker")},
			Owner:    owner("metadata"),
			Expected: 99,
		},
	)

	requireLifecycleError(t, err, ErrStaleRevision, ErrorReasonStaleRevision)
}

func updateObservedForPatchTest(t *testing.T, executor *Executor, expected objectstore.Revision) Result {
	t.Helper()

	result, err := executor.UpdateObserved(
		context.Background(),
		UpdateObservedRequest{
			Resource: testGVR(),
			Object:   testName(1),
			Observed: objectValue(member("ready", value.StringValue("true"))),
			Owner:    owner("controller"),
			Expected: expected,
		},
	)
	requireNoError(t, err)

	return result
}

func stringPtr(value string) *string {
	return &value
}

func requireLabel(t *testing.T, result Result, key string, want string) {
	t.Helper()

	got, ok := result.State.Object.ObjectMeta.Labels.Get(labels.Key(key))
	if !ok || got != labels.Value(want) {
		t.Fatalf("label %q = %q, %v; want %q, true", key, got, ok, want)
	}
}

func requireAnnotation(t *testing.T, result Result, key string, want string) {
	t.Helper()

	got, ok := result.State.Object.ObjectMeta.Annotations.Get(annotations.Key(key))
	if !ok || got != annotations.Value(want) {
		t.Fatalf("annotation %q = %q, %v; want %q, true", key, got, ok, want)
	}
}
