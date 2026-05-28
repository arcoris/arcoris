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
	var expr TypeExpr
	typ := ListOf(expr).Type()

	requireErrorIs(t, ValidateType(typ, nil), ErrInvalidType)
}

func TestListLengthAndSemantics(t *testing.T) {
	atomic := ListOf(String()).MinLen(1).MaxLen(3).Atomic().Type()
	set := ListOf(String()).Set().Type()

	requireNoError(t, ValidateType(atomic, nil))
	requireNoError(t, ValidateType(set, nil))

	view, ok := set.List()
	requireEqual(t, ok, true)
	requireEqual(t, view.Semantics(), ListSet)
}

func TestListInvalidLengthAndSemanticsRejected(t *testing.T) {
	invalidLen := ListOf(String()).MinLen(2).MaxLen(1).Type()
	invalidSemantics := ListOf(String()).Type()
	invalidSemantics.list.semantics = ListSemantics(99)

	requireErrorIs(t, ValidateType(invalidLen, nil), ErrInvalidType)
	requireErrorIs(t, ValidateType(invalidSemantics, nil), ErrInvalidType)
}

func TestListMapRequiresKeys(t *testing.T) {
	typ := ListOf(Object(Field("name").String().Required())).Map().Type()

	requireErrorIs(t, ValidateType(typ, nil), ErrInvalidField)
}

func TestListMapDirectObjectKeyValidation(t *testing.T) {
	valid := ListOf(Object(
		Field("type").String().Required(),
		Field("message").String().Optional(),
	)).Map("type").Type()
	missing := ListOf(Object(Field("type").String().Required())).Map("missing").Type()
	optional := ListOf(Object(Field("type").String().Optional())).Map("type").Type()

	requireNoError(t, ValidateType(valid, nil))
	requireErrorIs(t, ValidateType(missing, nil), ErrInvalidField)
	requireErrorIs(t, ValidateType(optional, nil), ErrInvalidField)
}

func TestListMapRefObjectKeyValidationWithResolver(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (TypeDefinition, bool) {
		if name == "example.Item" {
			return Define("example.Item", Object(
				Field("type").String().Required(),
				Field("value").String().Optional(),
			)), true
		}
		return TypeDefinition{}, false
	})

	typ := ListOf(Ref("example.Item")).Map("type").Type()
	requireNoError(t, ValidateType(typ, resolver))
}

func TestListMapKeysDetached(t *testing.T) {
	typ := ListOf(Object(Field("type").String().Required())).Map("type").Type()
	view, ok := typ.List()
	requireEqual(t, ok, true)
	keys := view.MapKeys()
	keys[0] = "changed"
	requireEqual(t, view.MapKeys()[0], FieldName("type"))
}

func TestListTypeExprMarker(t *testing.T) {
	ListOf(String()).typeExpr()
}
