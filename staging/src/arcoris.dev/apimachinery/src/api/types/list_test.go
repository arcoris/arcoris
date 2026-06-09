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

func TestListOfRequiresValidElement(t *testing.T) {
	var expr DescriptorExpr
	desc := ListOf(expr).Descriptor()

	requireErrorIs(t, ValidateLocal(desc), ErrInvalidDescriptor)
}

func TestListOfDefaultsToAtomic(t *testing.T) {
	desc := ListOf(String()).Descriptor()
	view, ok := desc.AsList()

	requireEqual(t, ok, true)
	requireEqual(t, view.Semantics(), ListAtomic)
	requireEqual(t, len(view.MapKeys()), 0)
}

func TestListTypeOrdered(t *testing.T) {
	desc := ListOf(String()).Ordered().Descriptor()
	view, ok := desc.AsList()

	requireEqual(t, ok, true)
	requireEqual(t, view.Semantics(), ListOrdered)
	requireEqual(t, len(view.MapKeys()), 0)
	requireNoError(t, ValidateLocal(desc))
}

func TestListOrderedClearsMapKeys(t *testing.T) {
	desc := ListOf(Object(Field("name").String().Required())).
		Map("name").
		Ordered().
		Descriptor()
	view, ok := desc.AsList()

	requireEqual(t, ok, true)
	requireEqual(t, view.Semantics(), ListOrdered)
	requireEqual(t, len(view.MapKeys()), 0)
	requireNoError(t, ValidateLocal(desc))
}

func TestListLengthAndSemantics(t *testing.T) {
	atomic := ListOf(String()).MinItems(1).MaxItems(3).Atomic().Descriptor()
	ordered := ListOf(String()).Ordered().Descriptor()
	set := ListOf(String()).Set().Descriptor()

	requireNoError(t, ValidateLocal(atomic))
	requireNoError(t, ValidateLocal(ordered))
	requireNoError(t, ValidateLocal(set))

	view, ok := set.AsList()
	requireEqual(t, ok, true)
	requireEqual(t, view.Semantics(), ListSet)
}

func TestListInvalidLengthAndSemanticsRejected(t *testing.T) {
	invalidLen := ListOf(String()).MinItems(2).MaxItems(1).Descriptor()
	invalidSemantics := ListOf(String()).Descriptor()
	invalidSemantics.list.semantics = ListSemantics(99)

	requireErrorIs(t, ValidateLocal(invalidLen), ErrInvalidDescriptor)
	requireErrorIs(t, ValidateLocal(invalidSemantics), ErrInvalidDescriptor)
}

func TestListMapRequiresKeys(t *testing.T) {
	desc := ListOf(Object(Field("name").String().Required())).Map().Descriptor()

	requireErrorIs(t, ValidateLocal(desc), ErrInvalidField)
}

func TestListMapDirectObjectKeyValidation(t *testing.T) {
	valid := ListOf(Object(
		Field("type").String().Required(),
		Field("message").String().Optional(),
	)).Map("type").Descriptor()
	missing := ListOf(Object(Field("type").String().Required())).Map("missing").Descriptor()
	optional := ListOf(Object(Field("type").String().Optional())).Map("type").Descriptor()

	requireNoError(t, ValidateLocal(valid))
	requireErrorIs(t, ValidateLocal(missing), ErrInvalidField)
	requireErrorIs(t, ValidateLocal(optional), ErrInvalidField)
}

func TestListMapRefObjectKeyValidationWithResolver(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		if name == "example.Item" {
			return Define("example.Item", Object(
				Field("type").String().Required(),
				Field("value").String().Optional(),
			)), true
		}
		return Definition{}, false
	})

	desc := ListOf(Ref("example.Item")).Map("type").Descriptor()
	requireNoError(t, ValidateResolved(desc, resolver))
}

func TestListMapKeysDetached(t *testing.T) {
	desc := ListOf(Object(Field("type").String().Required())).Map("type").Descriptor()
	view, ok := desc.AsList()
	requireEqual(t, ok, true)
	keys := view.MapKeys()
	keys[0] = "changed"
	requireEqual(t, view.MapKeys()[0], FieldName("type"))
}

func TestListDescriptorExprMarker(t *testing.T) {
	ListOf(String()).descriptorExpr()
}
