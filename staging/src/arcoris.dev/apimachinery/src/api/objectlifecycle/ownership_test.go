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

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuefieldset"
)

func TestCreateOwnershipInitFailureUsesPreciseReason(t *testing.T) {
	definition := resource.NewDefinition(
		testGroup,
		"Worker",
		"workers",
		resource.ScopeNamespaced,
		resource.NewVersion(
			"v1",
			types.Object(types.Field("image").Ref("example.dev.Image").Optional()).Descriptor(),
			resource.Exposed(),
			resource.Canonical(),
		),
	)
	executor, err := NewExecutor(
		WithStore(testStore(t)),
		WithResourceResolver(singleResourceResolver{definition: definition}),
		WithDesiredValidator(acceptingValueSurfaceValidator{}),
	)
	requireNoError(t, err)

	_, err = executor.Create(
		context.Background(),
		CreateRequest{Object: testObject(1, "api:v1"), Owner: owner("creator")},
	)

	requireLifecycleError(t, err, ErrApplyFailed, ErrorReasonOwnershipInitFailed)
	requireErrorIs(t, err, valuefieldset.ErrUnresolvedRef)
}

type acceptingValueSurfaceValidator struct{}

func (acceptingValueSurfaceValidator) ValidateSurface(value.Value, types.Descriptor, types.Resolver) error {
	return nil
}

type singleResourceResolver struct {
	definition resource.Definition
}

func (r singleResourceResolver) ResolveResource(gr apiidentity.GroupResource) (resource.Definition, bool) {
	if gr != r.definition.GroupResource() {
		return resource.Definition{}, false
	}

	return r.definition, true
}

func (r singleResourceResolver) ResolveKind(gk apiidentity.GroupKind) (resource.Definition, bool) {
	if gk != r.definition.GroupKind() {
		return resource.Definition{}, false
	}

	return r.definition, true
}

func (r singleResourceResolver) ResolveVersionResource(
	gvr apiidentity.GroupVersionResource,
) (resource.Definition, resource.VersionDefinition, bool) {
	if gvr.Group != r.definition.Group() || gvr.Resource != r.definition.Resource() {
		return resource.Definition{}, resource.VersionDefinition{}, false
	}
	version, ok := r.definition.Version(gvr.Version)
	if !ok {
		return resource.Definition{}, resource.VersionDefinition{}, false
	}

	return r.definition, version, true
}

func (r singleResourceResolver) ResolveVersionKind(
	gvk apiidentity.GroupVersionKind,
) (resource.Definition, resource.VersionDefinition, bool) {
	if gvk.Group != r.definition.Group() || gvk.Kind != r.definition.Kind() {
		return resource.Definition{}, resource.VersionDefinition{}, false
	}
	version, ok := r.definition.Version(gvk.Version)
	if !ok {
		return resource.Definition{}, resource.VersionDefinition{}, false
	}

	return r.definition, version, true
}
