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
	requireNoError(t, ValidateLocal(Ref("example.Name").Descriptor()))
}

func TestResolverValidationRejectsUnresolvedRef(t *testing.T) {
	resolver := resolverFunc(func(TypeName) (Definition, bool) {
		return Definition{}, false
	})

	requireErrorIs(t, ValidateResolved(Ref("example.Name").Descriptor(), resolver), ErrUnresolvedDescriptorReference)
}

func TestResolverValidationAcceptsResolvedRef(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		if name == "example.Name" {
			return Define("example.Name", String().MinBytes(1)), true
		}
		return Definition{}, false
	})

	requireNoError(t, ValidateResolved(Ref("example.Name").Descriptor(), resolver))
}

func TestResolverValidationRejectsResolvedInvalidDefinition(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		if name == "example.Bad" {
			return Define("example.Bad", ListOf(DescriptorExpr(nil))), true
		}
		return Definition{}, false
	})

	requireErrorIs(t, ValidateResolved(Ref("example.Bad").Descriptor(), resolver), ErrInvalidDescriptor)
}

func TestResolverValidationRejectsReferenceCycle(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		switch name {
		case "example.A":
			return Define("example.A", Ref("example.B")), true
		case "example.B":
			return Define("example.B", Ref("example.A")), true
		default:
			return Definition{}, false
		}
	})

	requireErrorIs(t, ValidateResolved(Ref("example.A").Descriptor(), resolver), ErrInvalidDescriptorReference)
}

func TestResolverValidationRejectsDirectDefinitionCycle(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		if name == "example.A" {
			return Define("example.A", Ref("example.A")), true
		}
		return Definition{}, false
	})

	requireErrorIs(t, ValidateDefinitionResolved(Define("example.A", Ref("example.A")), resolver), ErrInvalidDescriptorReference)
}

func TestResolverValidationSiblingRefsDoNotShareResolvingState(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		switch name {
		case "example.A":
			return Define("example.A", String().MinBytes(1)), true
		case "example.B":
			return Define("example.B", String().MinBytes(1)), true
		default:
			return Definition{}, false
		}
	})

	desc := Object(
		Field("a").Ref("example.A").Required(),
		Field("b").Ref("example.B").Required(),
	).Descriptor()

	requireValidDescriptor(t, desc, resolver)
}
