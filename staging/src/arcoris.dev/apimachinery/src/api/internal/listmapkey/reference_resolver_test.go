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

package listmapkey

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/internal/typeref"
	"arcoris.dev/apimachinery/api/types"
)

func TestReferenceResolverReturnsNonRefDescriptor(t *testing.T) {
	resolvedDescriptor, err := newReferenceResolver(Options{}).resolve(
		conditionPath(0),
		types.String().Type(),
		0,
	)

	requireNoError(t, err)
	requireEqual(t, resolvedDescriptor.Code(), types.TypeString)
}

func TestReferenceResolverResolvesReference(t *testing.T) {
	typeResolver := resolverFunc(func(name types.TypeName) (types.TypeDefinition, bool) {
		if name == "example.Name" {
			return types.Define("example.Name", types.String()), true
		}

		return types.TypeDefinition{}, false
	})

	resolvedDescriptor, err := newReferenceResolver(Options{Resolver: typeResolver}).resolve(
		conditionPath(0),
		types.Ref("example.Name").Type(),
		0,
	)

	requireNoError(t, err)
	requireEqual(t, resolvedDescriptor.Code(), types.TypeString)
}

func TestReferenceResolverResolvesReferenceChain(t *testing.T) {
	typeResolver := resolverFunc(func(name types.TypeName) (types.TypeDefinition, bool) {
		switch name {
		case "example.Name":
			return types.Define("example.Name", types.Ref("example.Text")), true
		case "example.Text":
			return types.Define("example.Text", types.String()), true
		default:
			return types.TypeDefinition{}, false
		}
	})

	resolvedDescriptor, err := newReferenceResolver(Options{Resolver: typeResolver}).resolve(
		conditionPath(0),
		types.Ref("example.Name").Type(),
		0,
	)

	requireNoError(t, err)
	requireEqual(t, resolvedDescriptor.Code(), types.TypeString)
}

func TestReferenceResolverRejectsMissingResolver(t *testing.T) {
	_, err := newReferenceResolver(Options{}).resolve(
		conditionPath(0),
		types.Ref("example.Name").Type(),
		0,
	)

	requireErrorKind(t, err, FailureUnresolvedRef)
}

func TestReferenceResolverRejectsResolverMiss(t *testing.T) {
	typeResolver := resolverFunc(func(types.TypeName) (types.TypeDefinition, bool) {
		return types.TypeDefinition{}, false
	})

	_, err := newReferenceResolver(Options{Resolver: typeResolver}).resolve(
		conditionPath(0),
		types.Ref("example.Name").Type(),
		0,
	)

	requireErrorKind(t, err, FailureUnresolvedRef)
}

func TestReferenceResolverRejectsReferenceCycle(t *testing.T) {
	typeResolver := resolverFunc(func(name types.TypeName) (types.TypeDefinition, bool) {
		switch name {
		case "example.A":
			return types.Define("example.A", types.Ref("example.B")), true
		case "example.B":
			return types.Define("example.B", types.Ref("example.A")), true
		default:
			return types.TypeDefinition{}, false
		}
	})

	_, err := newReferenceResolver(Options{Resolver: typeResolver}).resolve(
		conditionPath(0),
		types.Ref("example.A").Type(),
		0,
	)

	requireErrorKind(t, err, FailureReferenceCycle)
}

func TestReferenceResolverRejectsMaxDepth(t *testing.T) {
	typeResolver := resolverFunc(func(name types.TypeName) (types.TypeDefinition, bool) {
		return types.Define(name, types.String()), true
	})

	_, err := newReferenceResolver(Options{
		Resolver: typeResolver,
		MaxDepth: 1,
	}).resolve(
		conditionPath(0),
		types.Ref("example.Name").Type(),
		1,
	)

	requireErrorKind(t, err, FailureReferenceCycle)
}

func TestListMapRefErrorPreservesOtherErrors(t *testing.T) {
	want := errors.New("other")

	got := listMapRefError(want)

	if got != want {
		t.Fatalf("listMapRefError() = %v; want original error", got)
	}
}

func TestListMapRefErrorMapsTypeRefFailures(t *testing.T) {
	tests := []struct {
		name string
		kind typeref.FailureKind
		want FailureKind
	}{
		{name: "invalid descriptor", kind: typeref.FailureInvalidDescriptor, want: FailureInvalidDescriptor},
		{name: "unresolved ref", kind: typeref.FailureUnresolvedRef, want: FailureUnresolvedRef},
		{name: "reference cycle", kind: typeref.FailureReferenceCycle, want: FailureReferenceCycle},
		{name: "unknown kind", kind: typeref.FailureKind("other"), want: FailureInvalidDescriptor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := listMapRefError(&typeref.Error{
				Path:   conditionPath(0),
				Kind:   tt.kind,
				Detail: "detail",
			})

			requireErrorKind(t, err, tt.want)
		})
	}
}
