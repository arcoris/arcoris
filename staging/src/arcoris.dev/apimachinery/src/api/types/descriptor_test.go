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

func TestDescriptorZeroValueInvalid(t *testing.T) {
	var desc Descriptor

	requireEqual(t, desc.IsZero(), true)
	requireEqual(t, desc.IsValid(), false)
	requireEqual(t, desc.Code(), DescriptorInvalid)
	requireErrorIs(t, ValidateLocal(desc), ErrInvalidDescriptorKind)
}

func TestConstructorsCreateExpectedDescriptorKinds(t *testing.T) {
	tests := []struct {
		name       string
		descriptor Descriptor
		code       DescriptorKind
	}{
		{"null", Null().Descriptor(), DescriptorNull},
		{"bool", Bool().Descriptor(), DescriptorBool},
		{"string", String().Descriptor(), DescriptorString},
		{"bytes", Bytes().Descriptor(), DescriptorBytes},
		{"int8", Int8().Descriptor(), DescriptorInt8},
		{"int16", Int16().Descriptor(), DescriptorInt16},
		{"int32", Int32().Descriptor(), DescriptorInt32},
		{"int64", Int64().Descriptor(), DescriptorInt64},
		{"uint8", Uint8().Descriptor(), DescriptorUint8},
		{"uint16", Uint16().Descriptor(), DescriptorUint16},
		{"uint32", Uint32().Descriptor(), DescriptorUint32},
		{"uint64", Uint64().Descriptor(), DescriptorUint64},
		{"float32", Float32().Descriptor(), DescriptorFloat32},
		{"float64", Float64().Descriptor(), DescriptorFloat64},
		{"decimal", Decimal().Descriptor(), DescriptorDecimal},
		{"timestamp", Timestamp().Descriptor(), DescriptorTimestamp},
		{"date", Date().Descriptor(), DescriptorDate},
		{"time", Time().Descriptor(), DescriptorTime},
		{"duration", Duration().Descriptor(), DescriptorDuration},
		{"object", Object().Descriptor(), DescriptorObject},
		{"list", ListOf(String()).Descriptor(), DescriptorList},
		{"map", MapOf(String()).Descriptor(), DescriptorMap},
		{"ref", Ref("meta.arcoris.dev.Name").Descriptor(), DescriptorRef},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, tt.descriptor.Code(), tt.code)
			requireEqual(t, tt.descriptor.IsValid(), true)
		})
	}
}

func TestNullableFlagBehavior(t *testing.T) {
	tests := []Descriptor{
		Bool().Nullable().Descriptor(),
		String().Nullable().Descriptor(),
		Int64().Nullable().Descriptor(),
		Uint64().Nullable().Descriptor(),
		Float64().Nullable().Descriptor(),
		Decimal().Nullable().Descriptor(),
		Timestamp().Nullable().Descriptor(),
		Object().Nullable().Descriptor(),
		ListOf(String()).Nullable().Descriptor(),
		MapOf(String()).Nullable().Descriptor(),
		Ref("meta.arcoris.dev.Name").Nullable().Descriptor(),
	}
	for _, desc := range tests {
		requireEqual(t, desc.Nullable(), true)
	}
	requireEqual(t, Null().Descriptor().Nullable(), false)
}

func TestDescriptorNullCannotBeNullable(t *testing.T) {
	desc := Null().Descriptor()
	desc.flags = descriptorFlagNullable

	requireErrorIs(t, ValidateLocal(desc), ErrInvalidDescriptor)
}
