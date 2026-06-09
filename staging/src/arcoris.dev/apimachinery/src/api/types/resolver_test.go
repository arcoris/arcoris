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
type resolverFunc func(TypeName) (Definition, bool)

// Resolve implements Resolver.
func (f resolverFunc) Resolve(name TypeName) (Definition, bool) {
	return f(name)
}

func TestResolverValidationNilResolverAcceptsValidRefName(t *testing.T) {
	requireNoError(t, ValidateLocal(Ref("example.dev.Name").Descriptor()))
}

func TestResolverValidationResolvedRequiresResolver(t *testing.T) {
	requireErrorIs(t, ValidateResolved(Ref("example.dev.Name").Descriptor(), nil), ErrInvalidDescriptorReference)
}

func TestResolverValidationRejectsUnresolvedRef(t *testing.T) {
	resolver := resolverFunc(func(TypeName) (Definition, bool) {
		return Definition{}, false
	})

	requireErrorIs(t, ValidateResolved(Ref("example.dev.Name").Descriptor(), resolver), ErrUnresolvedDescriptorReference)
}

func TestResolverValidationAcceptsResolvedRef(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		if name == "example.dev.Name" {
			return Define("example.dev.Name", String().MinBytes(1)), true
		}
		return Definition{}, false
	})

	requireNoError(t, ValidateResolved(Ref("example.dev.Name").Descriptor(), resolver))
}

func TestResolverValidationRejectsResolvedInvalidDefinition(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		if name == "example.dev.Bad" {
			return Define("example.dev.Bad", ListOf(DescriptorExpr(nil))), true
		}
		return Definition{}, false
	})

	requireErrorIs(t, ValidateResolved(Ref("example.dev.Bad").Descriptor(), resolver), ErrInvalidDescriptor)
}

func TestResolverValidationRejectsReferenceCycle(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		switch name {
		case "example.dev.A":
			return Define("example.dev.A", Ref("example.dev.B")), true
		case "example.dev.B":
			return Define("example.dev.B", Ref("example.dev.A")), true
		default:
			return Definition{}, false
		}
	})

	requireErrorIs(t, ValidateResolved(Ref("example.dev.A").Descriptor(), resolver), ErrInvalidDescriptorReference)
}

func TestResolverValidationRejectsDirectDefinitionCycle(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		if name == "example.dev.A" {
			return Define("example.dev.A", Ref("example.dev.A")), true
		}
		return Definition{}, false
	})

	requireErrorIs(t, ValidateDefinitionResolved(Define("example.dev.A", Ref("example.dev.A")), resolver), ErrInvalidDescriptorReference)
}

func TestResolverValidationSiblingRefsDoNotShareResolvingState(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		switch name {
		case "example.dev.A":
			return Define("example.dev.A", String().MinBytes(1)), true
		case "example.dev.B":
			return Define("example.dev.B", String().MinBytes(1)), true
		default:
			return Definition{}, false
		}
	})

	desc := Object(
		Field("a").Ref("example.dev.A").Required(),
		Field("b").Ref("example.dev.B").Required(),
	).Descriptor()

	requireValidDescriptor(t, desc, resolver)
}
