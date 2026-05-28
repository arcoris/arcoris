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

package types

import "testing"

// resolverFunc adapts a function to Resolver for focused validation tests.
type resolverFunc func(TypeName) (TypeDefinition, bool)

// ResolveType implements Resolver.
func (f resolverFunc) ResolveType(name TypeName) (TypeDefinition, bool) {
	return f(name)
}

func TestResolverValidationNilResolverAcceptsValidRefName(t *testing.T) {
	requireNoError(t, ValidateType(Ref("example.Name").Type(), nil))
}

func TestResolverValidationRejectsUnresolvedRef(t *testing.T) {
	resolver := resolverFunc(func(TypeName) (TypeDefinition, bool) {
		return TypeDefinition{}, false
	})

	requireErrorIs(t, ValidateType(Ref("example.Name").Type(), resolver), ErrUnknownTypeReference)
}

func TestResolverValidationAcceptsResolvedRef(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (TypeDefinition, bool) {
		if name == "example.Name" {
			return Define("example.Name", String().MinLen(1)), true
		}
		return TypeDefinition{}, false
	})

	requireNoError(t, ValidateType(Ref("example.Name").Type(), resolver))
}

func TestResolverValidationRejectsResolvedInvalidDefinition(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (TypeDefinition, bool) {
		if name == "example.Bad" {
			return Define("example.Bad", ListOf(TypeExpr(nil))), true
		}
		return TypeDefinition{}, false
	})

	requireErrorIs(t, ValidateType(Ref("example.Bad").Type(), resolver), ErrInvalidType)
}

func TestResolverValidationRejectsReferenceCycle(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (TypeDefinition, bool) {
		switch name {
		case "example.A":
			return Define("example.A", Ref("example.B")), true
		case "example.B":
			return Define("example.B", Ref("example.A")), true
		default:
			return TypeDefinition{}, false
		}
	})

	requireErrorIs(t, ValidateType(Ref("example.A").Type(), resolver), ErrInvalidTypeReference)
}

func TestResolverValidationRejectsDirectDefinitionCycle(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (TypeDefinition, bool) {
		if name == "example.A" {
			return Define("example.A", Ref("example.A")), true
		}
		return TypeDefinition{}, false
	})

	requireErrorIs(t, ValidateDefinition(Define("example.A", Ref("example.A")), resolver), ErrInvalidTypeReference)
}

func TestResolverValidationSiblingRefsDoNotShareResolvingState(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (TypeDefinition, bool) {
		switch name {
		case "example.A":
			return Define("example.A", String().MinLen(1)), true
		case "example.B":
			return Define("example.B", String().MinLen(1)), true
		default:
			return TypeDefinition{}, false
		}
	})

	typ := Object(
		Field("a").Ref("example.A").Required(),
		Field("b").Ref("example.B").Required(),
	).Type()

	requireValidType(t, typ, resolver)
}
