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
		descriptor     Descriptor
		path           string
		reason         DescriptorErrorReason
		detailContains string
	}{
		{
			name:           "min non finite",
			descriptor:     Float64().Min(math.Inf(1)).Descriptor(),
			path:           "descriptor.min",
			reason:         DescriptorErrorReasonNonFiniteValue,
			detailContains: "+Inf",
		},
		{
			name:           "max non finite",
			descriptor:     Float64().Max(math.NaN()).Descriptor(),
			path:           "descriptor.max",
			reason:         DescriptorErrorReasonNonFiniteValue,
			detailContains: "NaN",
		},
		{
			name:           "range inverted",
			descriptor:     Float64().Range(10, 1).Descriptor(),
			path:           "descriptor.range",
			reason:         DescriptorErrorReasonInvalidRange,
			detailContains: "min=10 max=1",
		},
		{
			name:           "enum non finite",
			descriptor:     Float64().Enum(1, math.Inf(-1)).Descriptor(),
			path:           "descriptor.enum[1]",
			reason:         DescriptorErrorReasonNonFiniteValue,
			detailContains: "-Inf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireDescriptorError(t, ValidateLocal(tt.descriptor), ErrInvalidDescriptor, tt.path, tt.reason, tt.detailContains)
		})
	}
}

func TestValidationDiagnosticsEnumRules(t *testing.T) {
	requireDescriptorError(
		t,
		ValidateLocal(Int8().Enum(1, 1).Descriptor()),
		ErrInvalidDescriptor,
		"descriptor.enum",
		DescriptorErrorReasonDuplicateEnum,
		"duplicate value 1",
	)
	requireDescriptorError(
		t,
		ValidateLocal(Int8().Min(0).Enum(-1).Descriptor()),
		ErrInvalidDescriptor,
		"descriptor.enum[0]",
		DescriptorErrorReasonEnumBelowMinimum,
		"below minimum 0",
	)
	requireDescriptorError(
		t,
		ValidateLocal(Uint64().Max(1).Enum(2).Descriptor()),
		ErrInvalidDescriptor,
		"descriptor.enum[0]",
		DescriptorErrorReasonEnumAboveMaximum,
		"above maximum 1",
	)
}

func TestValidationDiagnosticsString(t *testing.T) {
	requireDescriptorError(
		t,
		ValidateLocal(String().Pattern("[").Descriptor()),
		ErrInvalidDescriptor,
		"descriptor.pattern",
		DescriptorErrorReasonInvalidPattern,
		"not a valid regexp",
	)
	requireDescriptorError(
		t,
		ValidateLocal(String().Pattern("^a+$").Enum("bbb").Descriptor()),
		ErrInvalidDescriptor,
		"descriptor.enum[0]",
		DescriptorErrorReasonEnumPatternMismatch,
		"does not match pattern",
	)
}

func TestValidationDiagnosticsBytesAndDecimal(t *testing.T) {
	requireDescriptorError(
		t,
		ValidateLocal(Bytes().MinBytes(-1).Descriptor()),
		ErrInvalidDescriptor,
		"descriptor.bytes.min",
		DescriptorErrorReasonNegativeLimit,
		"got -1",
	)
	requireDescriptorError(
		t,
		ValidateLocal(Bytes().MaxBytes(-1).Descriptor()),
		ErrInvalidDescriptor,
		"descriptor.bytes.max",
		DescriptorErrorReasonNegativeLimit,
		"got -1",
	)
	requireDescriptorError(
		t,
		ValidateLocal(Decimal().Precision(0).Descriptor()),
		ErrInvalidDescriptor,
		"descriptor.precision",
		DescriptorErrorReasonInvalidPrecision,
		"greater than zero",
	)
	requireDescriptorError(
		t,
		ValidateLocal(Decimal().Precision(2).Scale(3).Descriptor()),
		ErrInvalidDescriptor,
		"descriptor.scale",
		DescriptorErrorReasonInvalidScale,
		"scale=3 precision=2",
	)
}

func TestValidationDiagnosticsObjectListMapRefAndPayload(t *testing.T) {
	requireDescriptorError(
		t,
		ValidateLocal(Object(
			Field("name").String().Required(),
			Field("name").Int64().Optional(),
		).Descriptor()),
		ErrDuplicateField,
		"descriptor.fields[name].name",
		DescriptorErrorReasonDuplicateFieldName,
		"name",
	)

	requireDescriptorError(
		t,
		ValidateLocal(ListOf(Object(Field("name").String().Required())).Map().Descriptor()),
		ErrInvalidField,
		"descriptor.mapKeys",
		DescriptorErrorReasonMissingListMapKey,
		"requires at least one key",
	)
	requireDescriptorError(
		t,
		ValidateLocal(ListOf(Object(Field("name").String().Required())).Map("missing").Descriptor()),
		ErrInvalidField,
		"descriptor.mapKeys[0]",
		DescriptorErrorReasonListMapKeyNotFound,
		"not present",
	)
	requireDescriptorError(
		t,
		ValidateLocal(ListOf(Object(Field("name").String().Optional())).Map("name").Descriptor()),
		ErrInvalidField,
		"descriptor.mapKeys[0]",
		DescriptorErrorReasonListMapKeyOptional,
		"must be required",
	)

	mapping := MapOf(String()).Descriptor()
	mapping.mapType.value = nil
	requireDescriptorError(
		t,
		ValidateLocal(mapping),
		ErrInvalidDescriptor,
		"descriptor.value",
		DescriptorErrorReasonMissingValue,
		"must have a value descriptor",
	)

	missing := resolverFunc(func(TypeName) (Definition, bool) {
		return Definition{}, false
	})
	requireDescriptorError(
		t,
		ValidateResolved(Ref("example.dev.Missing").Descriptor(), missing),
		ErrUnresolvedDescriptorReference,
		"descriptor",
		DescriptorErrorReasonUnknownReference,
		"example.dev.Missing",
	)

	invalidResolved := resolverFunc(func(name TypeName) (Definition, bool) {
		if name == "example.dev.Bad" {
			return Define("example.dev.Bad", ListOf(DescriptorExpr(nil))), true
		}
		return Definition{}, false
	})
	requireDescriptorError(
		t,
		ValidateResolved(Ref("example.dev.Bad").Descriptor(), invalidResolved),
		ErrInvalidDescriptor,
		"descriptor",
		DescriptorErrorReasonInvalidResolvedDefinition,
		"example.dev.Bad",
	)

	cycle := resolverFunc(func(name TypeName) (Definition, bool) {
		switch name {
		case "example.dev.A":
			return Define("example.dev.A", Ref("example.dev.B")), true
		case "example.dev.B":
			return Define("example.dev.B", Ref("example.dev.A")), true
		default:
			return Definition{}, false
		}
	})
	requireDescriptorError(
		t,
		ValidateResolved(Ref("example.dev.A").Descriptor(), cycle),
		ErrInvalidDescriptorReference,
		"descriptor",
		DescriptorErrorReasonReferenceCycle,
		"recursive",
	)

	inactive := String().Descriptor()
	inactive.int8.min = limit[int8]{value: 1, set: true}
	requireDescriptorError(
		t,
		ValidateLocal(inactive),
		ErrInvalidDescriptor,
		"descriptor.payload",
		DescriptorErrorReasonInactivePayload,
		"int8 payload is populated",
	)
}
