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
	requireErrorIs(t, ValidateType(ListOf(TypeExpr(nil)).Type(), nil), ErrInvalidType)
	requireErrorIs(t, ValidateType(ListOf(String()).MinLen(2).MaxLen(1).Type(), nil), ErrInvalidType)
	requireErrorIs(t, ValidateType(ListOf(String()).Map().Type(), nil), ErrInvalidField)
	requireErrorIs(t, ValidateType(ListOf(String()).Map("type").Type(), nil), ErrInvalidType)

	invalidSemantics := ListOf(String()).Type()
	invalidSemantics.list.semantics = ListSemantics(99)
	requireErrorIs(t, ValidateType(invalidSemantics, nil), ErrInvalidType)
}

func TestListValidateMapKeys(t *testing.T) {
	valid := ListOf(Object(Field("type").String().Required())).Map("type").Type()
	missing := ListOf(Object(Field("type").String().Required())).Map("missing").Type()
	optional := ListOf(Object(Field("type").String().Optional())).Map("type").Type()

	requireNoError(t, ValidateType(valid, nil))
	requireErrorIs(t, ValidateType(missing, nil), ErrInvalidField)
	requireErrorIs(t, ValidateType(optional, nil), ErrInvalidField)
}

func TestListValidateRefMapKeys(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (TypeDefinition, bool) {
		switch name {
		case "example.Item":
			return Define("example.Item", Object(Field("type").String().Required())), true
		case "example.Name":
			return Define("example.Name", String()), true
		default:
			return TypeDefinition{}, false
		}
	})

	requireNoError(t, ValidateType(ListOf(Ref("example.Item")).Map("type").Type(), resolver))
	requireErrorIs(t, ValidateType(ListOf(Ref("example.Item")).Map("type").Type(), nil), ErrInvalidType)
	requireErrorIs(t, ValidateType(ListOf(Ref("example.Name")).Map("type").Type(), resolver), ErrInvalidType)
}
