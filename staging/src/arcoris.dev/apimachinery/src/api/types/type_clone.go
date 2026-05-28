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
// Type is a value object, but several exact payload slots contain slices to
// preserve declaration order. Every public boundary that returns or stores a
// Type calls cloneType so later caller mutation cannot rewrite descriptors held
// by catalogs, fields, or views.
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
