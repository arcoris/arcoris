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

// Type is the normalized structural descriptor / IR for an API value.
//
// Builders are the construction API; Type is the normalized descriptor they
// return. Its public shape is intentionally small: callers can read it through
// methods and exact read-only views, but cannot fill payload slots or attach
// arbitrary Go behavior. This keeps future validators, codecs, schema
// exporters, and resource-definition systems working from the same portable IR
// rather than from process-local implementation details.
//
// The zero value is invalid. Valid descriptors must be created through package
// constructors such as String, Int64, Object, ListOf, MapOf, and Ref, or through
// field builders that eventually produce object fields.
type Type struct {
	// code selects the active exact payload slot.
	code TypeCode

	// flags stores descriptor-wide flags such as nullability.
	flags typeFlags

	// string stores TypeString length, pattern, and enum rules.
	string stringPayload
	// bytes stores TypeBytes length rules.
	bytes bytesPayload
	// int8 stores TypeInt8 rules.
	int8 int8Payload
	// int16 stores TypeInt16 rules.
	int16 int16Payload
	// int32 stores TypeInt32 rules.
	int32 int32Payload
	// int64 stores TypeInt64 rules.
	int64 int64Payload
	// uint8 stores TypeUint8 rules.
	uint8 uint8Payload
	// uint16 stores TypeUint16 rules.
	uint16 uint16Payload
	// uint32 stores TypeUint32 rules.
	uint32 uint32Payload
	// uint64 stores TypeUint64 rules.
	uint64 uint64Payload
	// float32 stores TypeFloat32 rules.
	float32 float32Payload
	// float64 stores TypeFloat64 rules.
	float64 float64Payload
	// decimal stores TypeDecimal precision and scale rules.
	decimal decimalPayload
	// timestamp stores TypeTimestamp rules.
	timestamp timestampPayload
	// date stores TypeDate rules.
	date datePayload
	// timeOfDay stores TypeTime rules.
	timeOfDay timePayload
	// duration stores TypeDuration rules.
	duration durationPayload
	// object stores fixed object fields and unknown-field policy.
	object objectPayload
	// list stores element type, length rules, and list semantics.
	list listPayload
	// mapType stores dynamic map key and value descriptors.
	mapType mapPayload
	// ref stores the resolver name targeted by TypeRef.
	ref refPayload
}

// Code returns the structural category of t.
func (t Type) Code() TypeCode {
	return t.code
}

// IsZero reports whether t is the invalid zero descriptor.
func (t Type) IsZero() bool {
	return t.code == TypeInvalid
}

// IsValid reports whether t has a valid structural category.
//
// IsValid only checks the TypeCode. Full descriptor validation, including exact
// payload-slot consistency, scalar limits, object fields, list/map payloads,
// and reference resolution, belongs to ValidateType.
func (t Type) IsValid() bool {
	return t.code.IsValid()
}

// Nullable reports whether t admits null in addition to its structural value.
//
// Nullability is type-level value admissibility. It is independent from field
// presence: a required nullable field must be present, but its value may be
// null. TypeNull is the null literal itself and must never be marked nullable.
func (t Type) Nullable() bool {
	return t.flags&typeFlagNullable != 0
}
