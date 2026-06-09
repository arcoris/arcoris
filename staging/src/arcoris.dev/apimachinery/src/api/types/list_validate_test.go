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

import (
	"errors"
	"testing"
)

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
	duplicate := ListOf(Object(Field("type").String().Required())).Map("type", "type").Descriptor()

	requireNoError(t, ValidateLocal(valid))
	requireErrorIs(t, ValidateLocal(missing), ErrInvalidField)
	requireErrorIs(t, ValidateLocal(optional), ErrInvalidField)

	err := ValidateLocal(duplicate)
	requireDescriptorError(
		t,
		err,
		ErrInvalidField,
		"descriptor.mapKeys[1]",
		DescriptorErrorReasonDuplicateListMapKey,
		"duplicated at indexes 0 and 1",
	)

	var descriptorErr *DescriptorError
	requireEqual(t, errors.As(err, &descriptorErr), true)
	requireEqual(t, descriptorErr.Reason, DescriptorErrorReasonDuplicateListMapKey)
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
			desc: ListOf(Ref("meta.arcoris.dev.Name")).Set().Descriptor(),
			resolver: resolverFunc(func(name TypeName) (Definition, bool) {
				if name == "meta.arcoris.dev.Name" {
					return Define("meta.arcoris.dev.Name", String()), true
				}
				return Definition{}, false
			}),
		},
		{name: "local ref", desc: ListOf(Ref("meta.arcoris.dev.Name")).Set().Descriptor()},
		{name: "nullable", desc: ListOf(String().Nullable()).Set().Descriptor(), wantTarget: ErrInvalidDescriptor},
		{name: "object", desc: ListOf(Object()).Set().Descriptor(), wantTarget: ErrInvalidDescriptor},
		{name: "list", desc: ListOf(ListOf(String())).Set().Descriptor(), wantTarget: ErrInvalidDescriptor},
		{name: "map", desc: ListOf(MapOf(String())).Set().Descriptor(), wantTarget: ErrInvalidDescriptor},
		{name: "float", desc: ListOf(Float64()).Set().Descriptor(), wantTarget: ErrInvalidDescriptor},
		{name: "decimal", desc: ListOf(Decimal()).Set().Descriptor(), wantTarget: ErrInvalidDescriptor},
		{name: "temporal", desc: ListOf(Timestamp()).Set().Descriptor(), wantTarget: ErrInvalidDescriptor},
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

func TestListValidateSetElementRefResolvedFailures(t *testing.T) {
	missing := resolverFunc(func(TypeName) (Definition, bool) {
		return Definition{}, false
	})
	objectResolver := resolverFunc(func(name TypeName) (Definition, bool) {
		if name == "meta.arcoris.dev.Name" {
			return Define("meta.arcoris.dev.Name", Object()), true
		}
		return Definition{}, false
	})

	desc := ListOf(Ref("meta.arcoris.dev.Name")).Set().Descriptor()

	requireErrorIs(t, ValidateResolved(desc, missing), ErrUnresolvedDescriptorReference)
	requireDescriptorError(
		t,
		ValidateResolved(desc, objectResolver),
		ErrInvalidDescriptor,
		"descriptor.elem",
		DescriptorErrorReasonInvalidListSetElement,
		"descriptor object",
	)
}

func TestListValidateRefMapKeys(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		switch name {
		case "meta.arcoris.dev.Condition":
			return Define("meta.arcoris.dev.Condition", Object(Field("type").String().Required())), true
		case "meta.arcoris.dev.Name":
			return Define("meta.arcoris.dev.Name", String()), true
		default:
			return Definition{}, false
		}
	})

	requireNoError(t, ValidateLocal(ListOf(Ref("meta.arcoris.dev.Condition")).Map("type").Descriptor()))
	requireNoError(t, ValidateResolved(ListOf(Ref("meta.arcoris.dev.Condition")).Map("type").Descriptor(), resolver))
	requireErrorIs(t, ValidateResolved(ListOf(Ref("meta.arcoris.dev.Name")).Map("type").Descriptor(), resolver), ErrInvalidDescriptor)
}
