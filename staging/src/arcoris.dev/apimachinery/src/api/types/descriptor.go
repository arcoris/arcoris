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

// Descriptor is the normalized structural descriptor / IR for an API value.
//
// Builders are the construction API; Descriptor is the normalized descriptor they
// return. Its public shape is intentionally small: callers can read it through
// methods and exact read-only views, but cannot fill payload slots or attach
// arbitrary Go behavior. This keeps future validators, codecs, schema
// exporters, and resource-definition systems working from the same portable IR
// rather than from process-local implementation details.
//
// The zero value is invalid. Valid descriptors must be created through package
// constructors such as String, Int64, Object, ListOf, MapOf, and Ref, or through
// field builders that eventually produce object fields.
type Descriptor struct {
	// code selects the active exact payload slot.
	code DescriptorKind

	// flags stores descriptor-wide flags such as nullability.
	flags descriptorFlags

	// string stores DescriptorString length, pattern, and enum rules.
	string stringPayload
	// bytes stores DescriptorBytes length rules.
	bytes bytesPayload
	// int8 stores DescriptorInt8 rules.
	int8 int8Payload
	// int16 stores DescriptorInt16 rules.
	int16 int16Payload
	// int32 stores DescriptorInt32 rules.
	int32 int32Payload
	// int64 stores DescriptorInt64 rules.
	int64 int64Payload
	// uint8 stores DescriptorUint8 rules.
	uint8 uint8Payload
	// uint16 stores DescriptorUint16 rules.
	uint16 uint16Payload
	// uint32 stores DescriptorUint32 rules.
	uint32 uint32Payload
	// uint64 stores DescriptorUint64 rules.
	uint64 uint64Payload
	// float32 stores DescriptorFloat32 rules.
	float32 float32Payload
	// float64 stores DescriptorFloat64 rules.
	float64 float64Payload
	// decimal stores DescriptorDecimal precision and scale rules.
	decimal decimalPayload
	// timestamp stores DescriptorTimestamp rules.
	timestamp timestampPayload
	// date stores DescriptorDate rules.
	date datePayload
	// timeOfDay stores DescriptorTime rules.
	timeOfDay timePayload
	// duration stores DescriptorDuration rules.
	duration durationPayload
	// object stores fixed object fields and unknown-field policy.
	object objectPayload
	// list stores element descriptor, length rules, and list semantics.
	list listPayload
	// mapType stores dynamic map key and value descriptors.
	mapType mapPayload
	// ref stores the resolver name targeted by DescriptorRef.
	ref refPayload
}

// Code returns the structural category of desc.
func (desc Descriptor) Code() DescriptorKind {
	return desc.code
}

// String returns diagnostic text for the descriptor kind.
//
// String is intentionally not a view accessor. Use AsString, AsObject, AsList,
// and the other As* methods to inspect exact descriptor payloads.
func (desc Descriptor) String() string {
	return desc.code.String()
}

// IsZero reports whether desc is the invalid zero descriptor.
func (desc Descriptor) IsZero() bool {
	return desc.code == DescriptorInvalid
}

// IsValid reports whether desc has a valid structural category.
//
// IsValid only checks the DescriptorKind. Full descriptor validation, including exact
// payload-slot consistency, scalar limits, object fields, list/map payloads,
// and reference resolution, belongs to ValidateResolved.
func (desc Descriptor) IsValid() bool {
	return desc.code.IsValid()
}

// Nullable reports whether desc admits null in addition to its structural value.
//
// Nullability is descriptor-level value admissibility. It is independent from field
// presence: a required nullable field must be present, but its value may be
// null. DescriptorNull is the null literal itself and must never be marked nullable.
func (desc Descriptor) Nullable() bool {
	return desc.flags&descriptorFlagNullable != 0
}
