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

package resource

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/types"
)

func TestValidateDefinitionAcceptsValidCases(t *testing.T) {
	requireNoError(t, ValidateDefinitionLocal(validDefinition()))
	requireNoError(t, validDefinition().ValidateLocal())

	multi := NewDefinition(
		identity.Group("control.arcoris.dev"),
		identity.Kind("Worker"),
		identity.Resource("workers"),
		ScopeGlobal,
		NewVersion(identity.Version("v1alpha1"), objectType(), Exposed()),
		NewVersion(identity.Version("v1"), objectType(), Exposed(), Canonical()),
	)
	requireNoError(t, ValidateDefinitionLocal(multi))
	requireNoError(t, multi.ValidateLocal())
}

func TestValidateDefinitionAcceptsResolvedObjectRefs(t *testing.T) {
	resolver := fakeResolver{
		types.TypeName("control.arcoris.dev.WorkerDesired"): types.Define(
			"control.arcoris.dev.WorkerDesired",
			types.Object(),
		),
		types.TypeName("control.arcoris.dev.WorkerObserved"): types.Define(
			"control.arcoris.dev.WorkerObserved",
			types.Object(),
		),
	}

	def := NewDefinition(
		identity.Group("control.arcoris.dev"),
		identity.Kind("Worker"),
		identity.Resource("workers"),
		ScopeNamespaced,
		NewVersion(
			identity.Version("v1"),
			refType("control.arcoris.dev.WorkerDesired"),
			Observed(refType("control.arcoris.dev.WorkerObserved")),
			Exposed(),
			Canonical(),
		),
	)

	requireNoError(t, ValidateDefinitionResolved(def, resolver))
	requireNoError(t, def.ValidateResolved(resolver))
}

func TestValidateDefinitionLocalAcceptsRootRefs(t *testing.T) {
	def := NewDefinition(
		identity.Group("control.arcoris.dev"),
		identity.Kind("Worker"),
		identity.Resource("workers"),
		ScopeNamespaced,
		NewVersion(
			identity.Version("v1"),
			refType("control.arcoris.dev.WorkerDesired"),
			Observed(refType("control.arcoris.dev.WorkerObserved")),
			Exposed(),
			Canonical(),
		),
	)

	requireNoError(t, ValidateDefinitionLocal(def))
	requireNoError(t, def.ValidateLocal())
}

func TestValidateDefinitionPreservesNestedIdentityErrors(t *testing.T) {
	def := NewDefinition(
		identity.Group("apps"),
		identity.Kind("Worker"),
		identity.Resource("workers"),
		ScopeNamespaced,
		validVersion(),
	)

	err := ValidateDefinitionLocal(def)
	requireResourceError(t, err, ErrInvalidDefinition, pathDefinitionGroup, ErrorReasonInvalidGroup)

	if !errors.Is(err, identity.ErrInvalidIdentifier) {
		t.Fatalf("expected nested identity error, got %v", err)
	}
}
