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

// cloneType detaches all slice-bearing payloads from caller-owned state.
//
// Type is a value object, but several payload families contain slices to
// preserve declaration order. Every public boundary that returns or stores a
// Type calls cloneType so later caller mutation cannot rewrite descriptors held
// by registries, fields, or views.
func cloneType(t Type) Type {
	t.string = cloneStringPayload(t.string)
	t.bytes = cloneBytesPayload(t.bytes)
	t.int8 = cloneInt8Payload(t.int8)
	t.int16 = cloneInt16Payload(t.int16)
	t.int32 = cloneInt32Payload(t.int32)
	t.int64 = cloneInt64Payload(t.int64)
	t.uint8 = cloneUint8Payload(t.uint8)
	t.uint16 = cloneUint16Payload(t.uint16)
	t.uint32 = cloneUint32Payload(t.uint32)
	t.uint64 = cloneUint64Payload(t.uint64)
	t.float32 = cloneFloat32Payload(t.float32)
	t.float64 = cloneFloat64Payload(t.float64)
	t.decimal = cloneDecimalPayload(t.decimal)
	t.timestamp = cloneTimestampPayload(t.timestamp)
	t.date = cloneDatePayload(t.date)
	t.timeOfDay = cloneTimePayload(t.timeOfDay)
	t.duration = cloneDurationPayload(t.duration)
	t.object = cloneObjectPayload(t.object)
	t.list = cloneListPayload(t.list)
	t.mapType = cloneMapPayload(t.mapType)
	t.ref = cloneRefPayload(t.ref)
	return t
}

// cloneField detaches the Type payload stored inside f.
func cloneField(f FieldDescriptor) FieldDescriptor {
	f.typ = cloneType(f.typ)
	return f
}

// cloneFields detaches an ordered field list.
func cloneFields(fields []FieldDescriptor) []FieldDescriptor {
	if fields == nil {
		return nil
	}
	out := make([]FieldDescriptor, len(fields))
	for i := range fields {
		out[i] = cloneField(fields[i])
	}
	return out
}

// cloneFieldNames detaches an ordered field-name list.
func cloneFieldNames(names []FieldName) []FieldName {
	if names == nil {
		return nil
	}
	return append([]FieldName(nil), names...)
}

// cloneStrings detaches a string slice while preserving nil.
func cloneStrings(values []string) []string {
	if values == nil {
		return nil
	}
	return append([]string(nil), values...)
}

// cloneInt8s detaches an int8 slice while preserving nil.
func cloneInt8s(values []int8) []int8 {
	if values == nil {
		return nil
	}
	return append([]int8(nil), values...)
}

// cloneInt16s detaches an int16 slice while preserving nil.
func cloneInt16s(values []int16) []int16 {
	if values == nil {
		return nil
	}
	return append([]int16(nil), values...)
}

// cloneInt32s detaches an int32 slice while preserving nil.
func cloneInt32s(values []int32) []int32 {
	if values == nil {
		return nil
	}
	return append([]int32(nil), values...)
}

// cloneInt64s detaches an int64 slice while preserving nil.
func cloneInt64s(values []int64) []int64 {
	if values == nil {
		return nil
	}
	return append([]int64(nil), values...)
}

// cloneUint8s detaches a uint8 slice while preserving nil.
func cloneUint8s(values []uint8) []uint8 {
	if values == nil {
		return nil
	}
	return append([]uint8(nil), values...)
}

// cloneUint16s detaches a uint16 slice while preserving nil.
func cloneUint16s(values []uint16) []uint16 {
	if values == nil {
		return nil
	}
	return append([]uint16(nil), values...)
}

// cloneUint32s detaches a uint32 slice while preserving nil.
func cloneUint32s(values []uint32) []uint32 {
	if values == nil {
		return nil
	}
	return append([]uint32(nil), values...)
}

// cloneUint64s detaches a uint64 slice while preserving nil.
func cloneUint64s(values []uint64) []uint64 {
	if values == nil {
		return nil
	}
	return append([]uint64(nil), values...)
}

// cloneFloat32s detaches a float32 slice while preserving nil.
func cloneFloat32s(values []float32) []float32 {
	if values == nil {
		return nil
	}
	return append([]float32(nil), values...)
}

// cloneFloat64s detaches a float64 slice while preserving nil.
func cloneFloat64s(values []float64) []float64 {
	if values == nil {
		return nil
	}
	return append([]float64(nil), values...)
}
