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

package objectvalidation

import (
	"errors"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	apiobject "arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/types"
)

// testDesired is the simple requested-state surface shared by tests.
type testDesired struct {
	Replicas int32
}

// testObserved mirrors a simple computed/read surface used by contract tests.
type testObserved struct {
	ReadyReplicas int32
}

// testResolver lets tests verify that Plan.Resolver is passed through unchanged.
type testResolver struct {
	name string
}

// Resolve satisfies types.Resolver without resolving any references.
func (r *testResolver) Resolve(types.TypeName) (types.Definition, bool) {
	return types.Definition{}, false
}

// spySurfaceValidator records every value passed through the validation layer.
type spySurfaceValidator[T any] struct {
	called    int
	value     T
	typ       types.Descriptor
	resolver  types.Resolver
	err       error
	callOrder *[]string
	name      string
}

// ValidateSurface records call data and returns the configured error.
func (v *spySurfaceValidator[T]) ValidateSurface(
	value T,
	descriptor types.Descriptor,
	resolver types.Resolver,
) error {
	v.called++
	v.value = value
	v.typ = descriptor
	v.resolver = resolver

	if v.callOrder != nil {
		*v.callOrder = append(*v.callOrder, v.name)
	}

	return v.err
}

// desiredDescriptor returns the canonical desired descriptor used by fixtures.
func desiredDescriptor() types.Descriptor {
	return types.Object(
		types.Field("replicas").Int32().Required(),
	).Descriptor()
}

// observedDescriptor returns the canonical observed descriptor used by fixtures.
func observedDescriptor() types.Descriptor {
	return types.Object(
		types.Field("readyReplicas").Int32().Required(),
	).Descriptor()
}

// alternateDescriptor gives tests a distinct descriptor for selected-version checks.
func alternateDescriptor() types.Descriptor {
	return types.Object(
		types.Field("image").String().Required(),
	).Descriptor()
}

// versionWithObserved builds a version contract that defines observed state.
func versionWithObserved(version apiidentity.Version, options ...resource.VersionOption) resource.VersionDefinition {
	opts := append([]resource.VersionOption{resource.Observed(observedDescriptor())}, options...)
	return resource.NewVersion(version, desiredDescriptor(), opts...)
}

// versionWithoutObserved builds a version contract with desired state only.
func versionWithoutObserved(version apiidentity.Version, options ...resource.VersionOption) resource.VersionDefinition {
	return resource.NewVersion(version, desiredDescriptor(), options...)
}

// resourceDefinition builds the standard Worker resource used by tests.
func resourceDefinition(scope resource.Scope, versions ...resource.VersionDefinition) resource.Definition {
	if len(versions) == 0 {
		versions = []resource.VersionDefinition{versionWithObserved("v1")}
	}

	return resource.NewDefinition(
		apiidentity.Group("control.arcoris.dev"),
		apiidentity.Kind("Worker"),
		apiidentity.Resource("workers"),
		scope,
		versions...,
	)
}

// mismatchedResourceDefinition differs only by API group for match tests.
func mismatchedResourceDefinition() resource.Definition {
	return resource.NewDefinition(
		apiidentity.Group("other.arcoris.dev"),
		apiidentity.Kind("Worker"),
		apiidentity.Resource("workers"),
		resource.ScopeNamespaced,
		versionWithObserved("v1"),
	)
}

// validObject returns a namespaced Worker with desired and observed payloads.
func validObject() apiobject.Object[testDesired, testObserved] {
	return apiobject.NewObserved(
		validTypeMeta("v1"),
		validObjectMeta("system"),
		testDesired{Replicas: 3},
		testObserved{ReadyReplicas: 2},
	)
}

// validObjectWithoutObserved returns a namespaced Worker with desired payload only.
func validObjectWithoutObserved() apiobject.Object[testDesired, testObserved] {
	return apiobject.New[testDesired, testObserved](
		validTypeMeta("v1"),
		validObjectMeta("system"),
		testDesired{Replicas: 3},
	)
}

// validTypeMeta returns the GVK expected by the standard Worker resource.
func validTypeMeta(version apiidentity.Version) meta.TypeMeta {
	return meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
		Group:   "control.arcoris.dev",
		Version: version,
		Kind:    "Worker",
	})
}

// validObjectMeta returns object metadata with optional namespace.
func validObjectMeta(namespace metaidentity.Namespace) meta.ObjectMeta {
	return meta.ObjectMeta{
		Name:      "worker",
		Namespace: namespace,
		UID:       "uid-1",
	}
}

// validPlan returns a complete plan backed by recording validators.
func validPlan() Plan[testDesired, testObserved] {
	plan, _, _ := validPlanWithSpies()
	return plan
}

// validPlanWithSpies returns a complete plan and its validators for assertions.
func validPlanWithSpies() (
	Plan[testDesired, testObserved],
	*spySurfaceValidator[testDesired],
	*spySurfaceValidator[testObserved],
) {
	desired := &spySurfaceValidator[testDesired]{name: "desired"}
	observed := &spySurfaceValidator[testObserved]{name: "observed"}

	return Plan[testDesired, testObserved]{
		Resource:          resourceDefinition(resource.ScopeNamespaced),
		Resolver:          &testResolver{name: "primary"},
		DesiredValidator:  desired,
		ObservedValidator: observed,
	}, desired, observed
}

// requireNoError fails the test when err is non-nil.
func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// requireErrorIs asserts that err preserves target through errors.Is.
func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

// requireErrorNotIs asserts that err does not report target through errors.Is.
func requireErrorNotIs(t *testing.T, err error, target error) {
	t.Helper()
	if errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = true", err, target)
	}
}

// requireValidationError asserts the structured objectvalidation error shape.
func requireValidationError(
	t *testing.T,
	err error,
	target error,
	path string,
	reason ErrorReason,
) *Error {
	t.Helper()
	requireErrorIs(t, err, target)

	var validationErr *Error
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected *Error, got %T", err)
	}

	if validationErr.Path != path {
		t.Fatalf("Error.Path = %q, want %q", validationErr.Path, path)
	}
	if validationErr.Reason != reason {
		t.Fatalf("Error.Reason = %q, want %q", validationErr.Reason, reason)
	}
	if validationErr.Detail == "" {
		t.Fatal("Error.Detail is empty")
	}

	return validationErr
}

// requireCallOrder compares validator call order without hiding the sequence.
func requireCallOrder(t *testing.T, got []string, want ...string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("validator call count = %d, want %d; calls = %#v, want %#v", len(got), len(want), got, want)
	}

	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("validator calls[%d] = %q, want %q; calls = %#v, want %#v", i, got[i], want[i], got, want)
		}
	}
}

// requireAlternateDescriptor asserts the selected version descriptor through
// public descriptor views instead of depending on private Descriptor layout
// equality.
func requireAlternateDescriptor(t *testing.T, got types.Descriptor) {
	t.Helper()

	if got.Code() != types.DescriptorObject {
		t.Fatalf("Descriptor.Code() = %s, want %s", got.Code(), types.DescriptorObject)
	}

	view, ok := got.AsObject()
	if !ok {
		t.Fatal("Descriptor.AsObject() returned ok=false")
	}

	fields := view.Fields()
	if len(fields) != 1 {
		t.Fatalf("len(Object.Fields()) = %d, want 1; fields = %#v", len(fields), fields)
	}

	field := fields[0]
	if field.Name() != "image" {
		t.Fatalf("field.Name() = %q, want %q", field.Name(), "image")
	}
	if !field.IsRequired() {
		t.Fatal("field.IsRequired() = false, want true")
	}
	if field.Descriptor().Code() != types.DescriptorString {
		t.Fatalf("field.Descriptor().Code() = %s, want %s", field.Descriptor().Code(), types.DescriptorString)
	}
}
