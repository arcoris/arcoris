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

package typeref

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func TestResolveReturnsDefinition(t *testing.T) {
	name, descriptor, err := New(exampleResolver(), 64).Resolve(
		rootPath(),
		refType("example.Name"),
		0,
	)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if name != "example.Name" {
		t.Fatalf("Resolve() name = %q; want %q", name, "example.Name")
	}
	if descriptor.Code() != types.TypeString {
		t.Fatalf("Resolve() descriptor = %s; want %s", descriptor.Code(), types.TypeString)
	}
}

func TestResolveRejectsNonReferenceDescriptor(t *testing.T) {
	_, _, err := New(exampleResolver(), 64).Resolve(
		rootPath(),
		types.String().Type(),
		0,
	)
	requireFailureKind(t, err, FailureInvalidDescriptor)
}

func TestResolveRejectsMissingResolver(t *testing.T) {
	_, _, err := New(nil, 64).Resolve(
		rootPath(),
		refType("example.Name"),
		0,
	)
	requireFailureKind(t, err, FailureUnresolvedRef)
}

func TestResolveRejectsResolverMiss(t *testing.T) {
	resolver := resolverFunc(func(types.TypeName) (types.TypeDefinition, bool) {
		return types.TypeDefinition{}, false
	})

	_, _, err := New(resolver, 64).Resolve(
		rootPath(),
		refType("example.Name"),
		0,
	)
	requireFailureKind(t, err, FailureUnresolvedRef)
}

func TestResolveRejectsActiveReference(t *testing.T) {
	references := New(resolverFunc(stringResolver), 64)
	leave := references.Enter("example.Name")
	defer leave()

	_, _, err := references.Resolve(
		rootPath(),
		refType("example.Name"),
		0,
	)
	requireFailureKind(t, err, FailureReferenceCycle)
}

func TestResolveRejectsMaxDepth(t *testing.T) {
	_, _, err := New(resolverFunc(stringResolver), 1).Resolve(
		rootPath(),
		refType("example.Name"),
		1,
	)
	requireFailureKind(t, err, FailureReferenceCycle)
}
