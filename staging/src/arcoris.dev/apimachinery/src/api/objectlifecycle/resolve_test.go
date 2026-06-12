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
)

func TestResolveObjectResourceReportsInvalidResourceContract(t *testing.T) {
	executor := testExecutor(t, WithResourceResolver(inconsistentResourceResolver{}))

	_, err := executor.Create(
		context.Background(),
		CreateRequest{Object: testObject(1, "api:v1"), Owner: owner("creator")},
	)

	requireLifecycleError(t, err, ErrValidationFailed, ErrorReasonInvalidResourceContract)
	requireErrorIs(t, err, ErrInvalidResourceContract)
}

type inconsistentResourceResolver struct{}

func (inconsistentResourceResolver) ResolveResource(apiidentity.GroupResource) (resource.Definition, bool) {
	return resource.Definition{}, false
}

func (inconsistentResourceResolver) ResolveKind(apiidentity.GroupKind) (resource.Definition, bool) {
	return resource.Definition{}, false
}

func (inconsistentResourceResolver) ResolveVersionResource(
	apiidentity.GroupVersionResource,
) (resource.Definition, resource.VersionDefinition, bool) {
	return resource.Definition{}, resource.VersionDefinition{}, false
}

func (inconsistentResourceResolver) ResolveVersionKind(
	apiidentity.GroupVersionKind,
) (resource.Definition, resource.VersionDefinition, bool) {
	definition := resource.NewDefinition(
		testGroup,
		"Worker",
		"workers",
		resource.ScopeNamespaced,
		resource.NewVersion("v2", desiredDescriptor(), resource.Exposed(), resource.Canonical()),
	)

	return definition, resource.NewVersion("v1", desiredDescriptor()), true
}
