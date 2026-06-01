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

func TestListMapKeyValidationAcceptsStableIdentityScalars(t *testing.T) {
	tests := []struct {
		name     string
		field    FieldExpr
		resolver Resolver
	}{
		{name: "bool", field: Field("key").Bool().Required()},
		{name: "string", field: Field("key").String().Required()},
		{name: "int8", field: Field("key").Int8().Required()},
		{name: "int16", field: Field("key").Int16().Required()},
		{name: "int32", field: Field("key").Int32().Required()},
		{name: "int64", field: Field("key").Int64().Required()},
		{name: "uint8", field: Field("key").Uint8().Required()},
		{name: "uint16", field: Field("key").Uint16().Required()},
		{name: "uint32", field: Field("key").Uint32().Required()},
		{name: "uint64", field: Field("key").Uint64().Required()},
		{
			name:  "string ref",
			field: Field("key").Ref("example.StringKey").Required(),
			resolver: resolverFunc(func(name TypeName) (TypeDefinition, bool) {
				if name == "example.StringKey" {
					return Define("example.StringKey", String().MinLen(1)), true
				}
				return TypeDefinition{}, false
			}),
		},
		{
			name:  "uint64 ref",
			field: Field("key").Ref("example.Uint64Key").Required(),
			resolver: resolverFunc(func(name TypeName) (TypeDefinition, bool) {
				if name == "example.Uint64Key" {
					return Define("example.Uint64Key", Uint64()), true
				}
				return TypeDefinition{}, false
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := listMapWithKeyField(tt.field)

			requireNoError(t, ValidateType(typ, tt.resolver))
		})
	}
}

func TestListMapKeyValidationRejectsUnsupportedIdentityTypes(t *testing.T) {
	tests := []struct {
		name     string
		field    FieldExpr
		resolver Resolver
		detail   string
	}{
		{name: "null", field: Field("key").Null().Required(), detail: "type null"},
		{name: "bytes", field: Field("key").Bytes().Required(), detail: "type bytes"},
		{name: "float32", field: Field("key").Float32().Required(), detail: "type float32"},
		{name: "float64", field: Field("key").Float64().Required(), detail: "type float64"},
		{name: "decimal", field: Field("key").Decimal().Required(), detail: "type decimal"},
		{name: "timestamp", field: Field("key").Timestamp().Required(), detail: "type timestamp"},
		{name: "date", field: Field("key").Date().Required(), detail: "type date"},
		{name: "time", field: Field("key").Time().Required(), detail: "type time"},
		{name: "duration", field: Field("key").Duration().Required(), detail: "type duration"},
		{
			name: "object",
			field: Field("key").Object(
				Field("part").String().Required(),
			).Required(),
			detail: "type object",
		},
		{
			name:   "list",
			field:  Field("key").ListOf(String()).Required(),
			detail: "type list",
		},
		{
			name:   "map",
			field:  Field("key").MapOf(String()).Required(),
			detail: "type map",
		},
		{
			name:  "object ref",
			field: Field("key").Ref("example.ObjectKey").Required(),
			resolver: resolverFunc(func(name TypeName) (TypeDefinition, bool) {
				if name == "example.ObjectKey" {
					return Define("example.ObjectKey", Object(
						Field("part").String().Required(),
					)), true
				}
				return TypeDefinition{}, false
			}),
			detail: "type object",
		},
		{
			name:  "bytes ref",
			field: Field("key").Ref("example.BytesKey").Required(),
			resolver: resolverFunc(func(name TypeName) (TypeDefinition, bool) {
				if name == "example.BytesKey" {
					return Define("example.BytesKey", Bytes()), true
				}
				return TypeDefinition{}, false
			}),
			detail: "type bytes",
		},
		{
			name:   "nullable string",
			field:  Field("key").String().Nullable().Required(),
			detail: "non-nullable",
		},
		{
			name:  "nullable ref",
			field: Field("key").Ref("example.StringKey").Nullable().Required(),
			resolver: resolverFunc(func(name TypeName) (TypeDefinition, bool) {
				if name == "example.StringKey" {
					return Define("example.StringKey", String()), true
				}
				return TypeDefinition{}, false
			}),
			detail: "non-nullable",
		},
		{
			name:  "ref to nullable string",
			field: Field("key").Ref("example.NullableStringKey").Required(),
			resolver: resolverFunc(func(name TypeName) (TypeDefinition, bool) {
				if name == "example.NullableStringKey" {
					return Define("example.NullableStringKey", String().Nullable()), true
				}
				return TypeDefinition{}, false
			}),
			detail: "non-nullable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateType(listMapWithKeyField(tt.field), tt.resolver)

			requireTypeError(
				t,
				err,
				ErrInvalidField,
				"type.mapKeys[0]",
				TypeErrorReasonInvalidListMapKeyType,
				tt.detail,
			)
		})
	}
}

func TestListMapKeyValidationRejectsUnresolvedKeyRefWithoutResolver(t *testing.T) {
	err := ValidateType(
		listMapWithKeyField(Field("key").Ref("example.Key").Required()),
		nil,
	)

	requireTypeError(
		t,
		err,
		ErrUnknownTypeReference,
		"type.mapKeys[0]",
		TypeErrorReasonUnknownReference,
		"without a resolver",
	)
}

func TestListMapKeyValidationRejectsReferenceCycle(t *testing.T) {
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
	field := Field("key").Ref("example.A").Required().Field()

	err := validateListMapKeyIdentityType(
		field,
		resolver,
		"type.mapKeys[0]",
		make(map[TypeName]bool),
	)

	requireTypeError(
		t,
		err,
		ErrInvalidTypeReference,
		"type.mapKeys[0]",
		TypeErrorReasonReferenceCycle,
		"recursive",
	)
}

func TestListMapKeyValidationRejectsInvalidZeroType(t *testing.T) {
	field := FieldDescriptor{
		name:     "key",
		presence: PresenceRequired,
		typ:      Type{},
	}

	err := validateListMapKeyIdentityType(
		field,
		nil,
		"type.mapKeys[0]",
		make(map[TypeName]bool),
	)

	requireTypeError(
		t,
		err,
		ErrInvalidField,
		"type.mapKeys[0]",
		TypeErrorReasonInvalidListMapKeyType,
		"type invalid",
	)
}

func listMapWithKeyField(field FieldExpr) Type {
	return ListOf(Object(
		field,
		Field("value").String().Required(),
	)).Map("key").Type()
}
