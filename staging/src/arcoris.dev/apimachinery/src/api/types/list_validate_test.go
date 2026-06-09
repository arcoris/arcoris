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

func TestListValidateRejectsInvalidShapes(t *testing.T) {
	requireErrorIs(t, ValidateLocal(ListOf(DescriptorExpr(nil)).Descriptor()), ErrInvalidDescriptor)
	requireErrorIs(t, ValidateLocal(ListOf(String()).MinItems(2).MaxItems(1).Descriptor()), ErrInvalidDescriptor)
	requireErrorIs(t, ValidateLocal(ListOf(String()).Map().Descriptor()), ErrInvalidField)
	requireErrorIs(t, ValidateLocal(ListOf(String()).Map("type").Descriptor()), ErrInvalidDescriptor)

	invalidSemantics := ListOf(String()).Descriptor()
	invalidSemantics.list.semantics = ListSemantics(99)
	requireErrorIs(t, ValidateLocal(invalidSemantics), ErrInvalidDescriptor)
}

func TestListValidateMapKeys(t *testing.T) {
	valid := ListOf(Object(Field("type").String().Required())).Map("type").Descriptor()
	missing := ListOf(Object(Field("type").String().Required())).Map("missing").Descriptor()
	optional := ListOf(Object(Field("type").String().Optional())).Map("type").Descriptor()

	requireNoError(t, ValidateLocal(valid))
	requireErrorIs(t, ValidateLocal(missing), ErrInvalidField)
	requireErrorIs(t, ValidateLocal(optional), ErrInvalidField)
}

func TestValidateListOrderedAcceptsNoMapKeys(t *testing.T) {
	desc := ListOf(String()).Ordered().Descriptor()

	requireNoError(t, ValidateLocal(desc))
}

func TestValidateListOrderedRejectsUnexpectedMapKeysIfConstructible(t *testing.T) {
	desc := ListOf(String()).Ordered().Descriptor()
	desc.list.mapKeys = []FieldName{"name"}

	requireErrorIs(t, ValidateLocal(desc), ErrInvalidField)
}

func TestListValidateSetElements(t *testing.T) {
	cases := []struct {
		name       string
		desc       Descriptor
		resolver   Resolver
		wantTarget error
	}{
		{name: "bool", desc: ListOf(Bool()).Set().Descriptor()},
		{name: "string", desc: ListOf(String()).Set().Descriptor()},
		{name: "int", desc: ListOf(Int64()).Set().Descriptor()},
		{name: "uint", desc: ListOf(Uint64()).Set().Descriptor()},
		{
			name: "ref stable scalar",
			desc: ListOf(Ref("example.Name")).Set().Descriptor(),
			resolver: resolverFunc(func(name TypeName) (Definition, bool) {
				if name == "example.Name" {
					return Define("example.Name", String()), true
				}
				return Definition{}, false
			}),
		},
		{name: "nullable", desc: ListOf(String().Nullable()).Set().Descriptor(), wantTarget: ErrInvalidDescriptor},
		{name: "object", desc: ListOf(Object()).Set().Descriptor(), wantTarget: ErrInvalidDescriptor},
		{name: "list", desc: ListOf(ListOf(String())).Set().Descriptor(), wantTarget: ErrInvalidDescriptor},
		{name: "map", desc: ListOf(MapOf(String())).Set().Descriptor(), wantTarget: ErrInvalidDescriptor},
		{name: "float", desc: ListOf(Float64()).Set().Descriptor(), wantTarget: ErrInvalidDescriptor},
		{name: "decimal", desc: ListOf(Decimal()).Set().Descriptor(), wantTarget: ErrInvalidDescriptor},
		{name: "temporal", desc: ListOf(Timestamp()).Set().Descriptor(), wantTarget: ErrInvalidDescriptor},
		{name: "unresolved ref", desc: ListOf(Ref("example.Name")).Set().Descriptor(), wantTarget: ErrUnresolvedDescriptorReference},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantTarget == nil {
				requireNoError(t, validateTestDescriptor(tc.desc, tc.resolver))
				return
			}

			requireErrorIs(t, validateTestDescriptor(tc.desc, tc.resolver), tc.wantTarget)
		})
	}
}

func TestListValidateRefMapKeys(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		switch name {
		case "example.Item":
			return Define("example.Item", Object(Field("type").String().Required())), true
		case "example.Name":
			return Define("example.Name", String()), true
		default:
			return Definition{}, false
		}
	})

	requireNoError(t, ValidateResolved(ListOf(Ref("example.Item")).Map("type").Descriptor(), resolver))
	requireErrorIs(t, ValidateLocal(ListOf(Ref("example.Item")).Map("type").Descriptor()), ErrInvalidDescriptor)
	requireErrorIs(t, ValidateResolved(ListOf(Ref("example.Name")).Map("type").Descriptor(), resolver), ErrInvalidDescriptor)
}
