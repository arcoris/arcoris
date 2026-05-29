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
	"math"
	"testing"
)

func TestValidationDiagnosticsFloat64(t *testing.T) {
	tests := []struct {
		name           string
		typ            Type
		path           string
		reason         TypeErrorReason
		detailContains string
	}{
		{
			name:           "min non finite",
			typ:            Float64().Min(math.Inf(1)).Type(),
			path:           "type.min",
			reason:         TypeErrorReasonNonFiniteValue,
			detailContains: "+Inf",
		},
		{
			name:           "max non finite",
			typ:            Float64().Max(math.NaN()).Type(),
			path:           "type.max",
			reason:         TypeErrorReasonNonFiniteValue,
			detailContains: "NaN",
		},
		{
			name:           "range inverted",
			typ:            Float64().Range(10, 1).Type(),
			path:           "type.range",
			reason:         TypeErrorReasonInvalidRange,
			detailContains: "min=10 max=1",
		},
		{
			name:           "enum non finite",
			typ:            Float64().Enum(1, math.Inf(-1)).Type(),
			path:           "type.enum[1]",
			reason:         TypeErrorReasonNonFiniteValue,
			detailContains: "-Inf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireTypeError(t, ValidateType(tt.typ, nil), ErrInvalidType, tt.path, tt.reason, tt.detailContains)
		})
	}
}

func TestValidationDiagnosticsEnumRules(t *testing.T) {
	requireTypeError(
		t,
		ValidateType(Int8().Enum(1, 1).Type(), nil),
		ErrInvalidType,
		"type.enum",
		TypeErrorReasonDuplicateEnum,
		"duplicate value 1",
	)
	requireTypeError(
		t,
		ValidateType(Int8().Min(0).Enum(-1).Type(), nil),
		ErrInvalidType,
		"type.enum[0]",
		TypeErrorReasonEnumBelowMinimum,
		"below minimum 0",
	)
	requireTypeError(
		t,
		ValidateType(Uint64().Max(1).Enum(2).Type(), nil),
		ErrInvalidType,
		"type.enum[0]",
		TypeErrorReasonEnumAboveMaximum,
		"above maximum 1",
	)
}

func TestValidationDiagnosticsString(t *testing.T) {
	requireTypeError(
		t,
		ValidateType(String().Pattern("[").Type(), nil),
		ErrInvalidType,
		"type.pattern",
		TypeErrorReasonInvalidPattern,
		"not a valid regexp",
	)
	requireTypeError(
		t,
		ValidateType(String().Pattern("^a+$").Enum("bbb").Type(), nil),
		ErrInvalidType,
		"type.enum[0]",
		TypeErrorReasonEnumPatternMismatch,
		"does not match pattern",
	)
}

func TestValidationDiagnosticsBytesAndDecimal(t *testing.T) {
	requireTypeError(
		t,
		ValidateType(Bytes().MinLen(-1).Type(), nil),
		ErrInvalidType,
		"type.len.min",
		TypeErrorReasonNegativeLimit,
		"got -1",
	)
	requireTypeError(
		t,
		ValidateType(Bytes().MaxLen(-1).Type(), nil),
		ErrInvalidType,
		"type.len.max",
		TypeErrorReasonNegativeLimit,
		"got -1",
	)
	requireTypeError(
		t,
		ValidateType(Decimal().Precision(0).Type(), nil),
		ErrInvalidType,
		"type.precision",
		TypeErrorReasonInvalidPrecision,
		"greater than zero",
	)
	requireTypeError(
		t,
		ValidateType(Decimal().Precision(2).Scale(3).Type(), nil),
		ErrInvalidType,
		"type.scale",
		TypeErrorReasonInvalidScale,
		"scale=3 precision=2",
	)
}

func TestValidationDiagnosticsObjectListMapRefAndPayload(t *testing.T) {
	requireTypeError(
		t,
		ValidateType(Object(
			Field("name").String().Required(),
			Field("name").Int64().Optional(),
		).Type(), nil),
		ErrDuplicateField,
		"type.fields[name].name",
		TypeErrorReasonDuplicateFieldName,
		"name",
	)

	requireTypeError(
		t,
		ValidateType(ListOf(Object(Field("name").String().Required())).Map().Type(), nil),
		ErrInvalidField,
		"type.mapKeys",
		TypeErrorReasonMissingListMapKey,
		"requires at least one key",
	)
	requireTypeError(
		t,
		ValidateType(ListOf(Object(Field("name").String().Required())).Map("missing").Type(), nil),
		ErrInvalidField,
		"type.mapKeys[0]",
		TypeErrorReasonListMapKeyNotFound,
		"not present",
	)
	requireTypeError(
		t,
		ValidateType(ListOf(Object(Field("name").String().Optional())).Map("name").Type(), nil),
		ErrInvalidField,
		"type.mapKeys[0]",
		TypeErrorReasonListMapKeyOptional,
		"must be required",
	)

	mapping := MapOf(String()).Type()
	mapping.mapType.value = nil
	requireTypeError(
		t,
		ValidateType(mapping, nil),
		ErrInvalidType,
		"type.value",
		TypeErrorReasonMissingValue,
		"must have a value type",
	)

	missing := resolverFunc(func(TypeName) (TypeDefinition, bool) {
		return TypeDefinition{}, false
	})
	requireTypeError(
		t,
		ValidateType(Ref("example.Missing").Type(), missing),
		ErrUnknownTypeReference,
		"type",
		TypeErrorReasonUnknownReference,
		"example.Missing",
	)

	invalidResolved := resolverFunc(func(name TypeName) (TypeDefinition, bool) {
		if name == "example.Bad" {
			return Define("example.Bad", ListOf(TypeExpr(nil))), true
		}
		return TypeDefinition{}, false
	})
	requireTypeError(
		t,
		ValidateType(Ref("example.Bad").Type(), invalidResolved),
		ErrInvalidType,
		"type",
		TypeErrorReasonInvalidResolvedDefinition,
		"example.Bad",
	)

	cycle := resolverFunc(func(name TypeName) (TypeDefinition, bool) {
		switch name {
		case "example.A":
			return Define("example.A", Ref("example.B")), true
		case "example.B":
			return Define("example.B", Ref("example.A")), true
		default:
			return TypeDefinition{}, false
		}
	})
	requireTypeError(
		t,
		ValidateType(Ref("example.A").Type(), cycle),
		ErrInvalidTypeReference,
		"type",
		TypeErrorReasonReferenceCycle,
		"recursive",
	)

	inactive := String().Type()
	inactive.int8.min = limit[int8]{value: 1, set: true}
	requireTypeError(
		t,
		ValidateType(inactive, nil),
		ErrInvalidType,
		"type.payload",
		TypeErrorReasonInactivePayload,
		"int8 payload is populated",
	)
}
